package services

import (
	"fmt"
	"libreria/graph/model"
	"libreria/models"
	"strconv"

	"github.com/couchbase/gocb/v2"
)

func TomarRecurso(id int) (model.Recurso, error) {
	var recurso model.Recurso

	// Obtener datos de MariaDB
	err := models.MariaDB.QueryRow(`
		SELECT * FROM documentos WHERE idDoc = ?`, id).Scan(
		&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria,
		&recurso.IDUsuario, &recurso.Formato, &recurso.Descripcion,
		&recurso.FechaOrigen, &recurso.DireccionRec,
	)
	if err != nil {
		return recurso, fmt.Errorf("error al obtener datos de MariaDB: %w", err)
	}

	// Construir el ID de documento
	docID := "documento::" + strconv.Itoa(id)

	// Consulta en Couchbase
	query := "SELECT etiquetas, calificacion, noDescargas FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return recurso, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}

	// Procesar resultado de Couchbase
	resultado, err := devolverResultadoClaveValor(rows)
	if err != nil {
		return recurso, fmt.Errorf("error al procesar el resultado de Couchbase: %w", err)
	}

	// Manejar etiquetas
	if etiquetasRaw, ok := resultado["etiquetas"].([]interface{}); ok {
		etiquetas, err := manejoEtiquetas(etiquetasRaw)
		if err != nil {
			return recurso, fmt.Errorf("error al procesar las etiquetas: %w", err)
		}
		recurso.Etiquetas = etiquetas
	}

	// Manejar calificación
	if calificacionRaw, ok := resultado["calificacion"].(float64); ok {
		recurso.Calificacion = &calificacionRaw
	}

	// Manejar descargas
	if noDescargasRaw, ok := resultado["noDescargas"].(float64); ok {
		noDescargas := int(noDescargasRaw) // Convertir float64 a int
		recurso.NumDescargas = &noDescargas
	}

	return recurso, nil
}

func BuscarRecursos(parametros model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	// Consulta en MariaDB
	rows, err := models.MariaDB.Query(`
		SELECT idDoc, titulo, autor, categoria, formato
		FROM documentos 
		WHERE 
			(titulo = ? OR ? IS NULL) AND
			(autor = ? OR ? IS NULL) AND
			(categoria = ? OR ? IS NULL) AND
			(formato = ? OR ? IS NULL);
	`, parametros.Titulo, parametros.Titulo, parametros.Autor, parametros.Autor, parametros.Categoria, parametros.Categoria, parametros.Formato, parametros.Formato)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en MariaDB: %w", err)
	}
	defer rows.Close()

	// Acumular resultados de MariaDB
	var recursos []*model.RecursoMuestra
	var ids []string

	for rows.Next() {
		var recurso model.RecursoMuestra
		if err := rows.Scan(&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato); err != nil {
			return nil, fmt.Errorf("error al leer los resultados de MariaDB: %w", err)
		}
		recursos = append(recursos, &recurso)
		ids = append(ids, "documento::"+strconv.Itoa(recurso.ID))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados de MariaDB: %w", err)
	}

	// Consulta en Couchbase con todos los IDs
	// Verificar si parametros.Etiquetas es nulo
	query := "SELECT META().id AS id, etiquetas, calificacion FROM `EduAndTime`.`documentos`.`documento` WHERE META().id IN $1"
	positionalParams := []interface{}{ids}

	// Agregar condición de etiquetas solo si parametros.Etiquetas no es nulo
	if parametros.Etiquetas != nil && len(parametros.Etiquetas) > 0 {
		query += " AND ANY etiqueta IN etiquetas SATISFIES etiqueta IN $2 END"
		positionalParams = append(positionalParams, parametros.Etiquetas)
	}

	documentos, err := models.CouchbaseCluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: positionalParams,
	})
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}

	// Procesar los resultados de Couchbase
	documentoMap := make(map[string]map[string]interface{})
	for documentos.Next() {
		var resultado map[string]interface{}
		if err := documentos.Row(&resultado); err != nil {
			return nil, fmt.Errorf("error al leer los resultados de Couchbase: %w", err)
		}
		id, ok := resultado["id"].(string)
		if !ok {
			return nil, fmt.Errorf("el campo 'id' no está presente o no es un string en el resultado de Couchbase")
		}
		documentoMap[id] = resultado
	}
	if err := documentos.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados de Couchbase: %w", err)
	}

	// Combinar resultados de MariaDB y Couchbase
	for _, recurso := range recursos {
		id := "documento::" + strconv.Itoa(recurso.ID)
		doc, exists := documentoMap[id]
		if !exists {
			continue
		}

		// Manejar etiquetas
		etiquetasRaw, ok := doc["etiquetas"].([]interface{})
		if ok {
			etiquetas, err := manejoEtiquetas(etiquetasRaw)
			if err != nil {
				return nil, fmt.Errorf("error al procesar etiquetas para el ID %s: %w", id, err)
			}
			recurso.Etiquetas = etiquetas
		}

		// Manejar calificación
		if calificacionRaw, ok := doc["calificacion"].(float64); ok {
			recurso.Calificacion = &calificacionRaw
		}
	}

	return recursos, nil
}

func UltimosRecursos() ([]*model.RecursoMuestra, error) {
	var recursos []*model.RecursoMuestra

	// Consulta para obtener los últimos 3 recursos
	rows, err := models.MariaDB.Query(`
		SELECT idDoc, titulo, autor, categoria, formato 
		FROM documentos 
		ORDER BY idDoc DESC 
		LIMIT 3;`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener datos de MariaDB: %w", err)
	}
	defer rows.Close()

	// Iterar sobre los resultados
	for rows.Next() {
		var recurso model.RecursoMuestra

		// Leer los datos del recurso desde MariaDB
		if err := rows.Scan(
			&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato,
		); err != nil {
			return nil, fmt.Errorf("error al leer los datos de MariaDB: %w", err)
		}

		// Completar los datos desde Couchbase
		completado, err := devolverDatosCouchbase(recurso)
		if err != nil {
			return nil, fmt.Errorf("error al obtener datos de Couchbase para el recurso %d: %w", recurso.ID, err)
		}

		// Agregar el recurso completado a la lista
		recursos = append(recursos, &completado)
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de los resultados: %w", err)
	}

	return recursos, nil
}

func RecursosAleatorios() ([]*model.RecursoMuestra, error) {
	var recursos []*model.RecursoMuestra

	// Consulta para obtener 3 recursos aleatorios
	rows, err := models.MariaDB.Query(`
		SELECT idDoc, titulo, autor, categoria, formato 
		FROM documentos 
		ORDER BY RAND() 
		LIMIT 3;`)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta aleatoria en MariaDB: %w", err)
	}
	defer rows.Close()

	// Iterar sobre los resultados
	for rows.Next() {
		var recurso model.RecursoMuestra

		// Leer los datos básicos del recurso desde MariaDB
		if err := rows.Scan(
			&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato,
		); err != nil {
			return nil, fmt.Errorf("error al leer datos de un recurso en MariaDB: %w", err)
		}

		// Completar datos desde Couchbase
		completado, err := devolverDatosCouchbase(recurso)
		if err != nil {
			return nil, fmt.Errorf("error al obtener datos de Couchbase para el recurso %d: %w", recurso.ID, err)
		}

		// Agregar el recurso completado a la lista
		recursos = append(recursos, &completado)
	}

	// Verificar errores durante la iteración
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración de resultados: %w", err)
	}

	return recursos, nil
}
