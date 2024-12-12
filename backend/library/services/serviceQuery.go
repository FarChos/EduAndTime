package services

import (
	"fmt"
	"libreria/graph/model"
	"libreria/models"
	"strconv"

	"github.com/couchbase/gocb/v2"
)

const TYPE_DOC = "documento::"

// !
func TomarDeMariadb(id int) (model.Recurso, error) {
	var recurso model.Recurso

	// Obtener datos de MariaDB
	err := models.MariaDB.QueryRow(`
		SELECT * FROM documentos WHERE idDoc = ?`, id).Scan(
		&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria,
		&recurso.IDUsuario, &recurso.Formato, &recurso.Descripcion,
		&recurso.FechaOrigen, &recurso.Archivo,
	)
	if err != nil {
		return recurso, fmt.Errorf("error al obtener datos de MariaDB: %w", err)
	}
	return recurso, nil
}

// ! TomarDeCouchbase
func TomarDeCouchbase(recurso model.Recurso) (model.Recurso, error) {
	docID := TYPE_DOC + strconv.Itoa(recurso.ID)

	query := "SELECT etiquetas, calificacion, noDescargas FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return recurso, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}
	defer rows.Close()

	// Procesar resultados
	for rows.Next() {
		var resultado map[string]interface{}
		if err := rows.Row(&resultado); err != nil {
			return recurso, fmt.Errorf("error al procesar la fila de resultados: %w", err)
		}

		// Manejar etiquetas
		if etiquetasRaw, ok := resultado["etiquetas"].([]interface{}); ok {
			etiquetas, err := ManejoEtiquetas(etiquetasRaw)
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
			noDescargas := int(noDescargasRaw)
			recurso.NumDescargas = &noDescargas
		}
	}

	if err := rows.Err(); err != nil {
		return recurso, fmt.Errorf("error al iterar sobre las filas: %w", err)
	}

	return recurso, nil
}

// !TomarNoDescargasEnCouchBase
func TomarNoDescargasEnCouchBase(id int) (*gocb.QueryResult, error) {
	docID := TYPE_DOC + strconv.FormatInt(int64(id), 10)

	// Consulta para obtener noDescargas
	queryRead := "SELECT noDescargas FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(queryRead, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en Couchbase : %w", err)
	}
	return rows, nil
}

// ! obtenerCalificaciones
func ObtenerCalificaciones(id int) (*gocb.QueryResult, error) {
	// Construir el ID del documento
	docID := TYPE_DOC + strconv.Itoa(id)

	// Consulta para obtener noCalificaciones y calificacion
	queryRead := "SELECT noCalificaciones, calificacion FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(queryRead, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}
	return rows, nil
}

// !
func TomarMuestraconIdes(ides []int) ([]*model.RecursoMuestra, error) {
	// Verificar si la lista de IDs está vacía
	if len(ides) == 0 {
		return nil, fmt.Errorf("la lista de IDs está vacía")
	}

	// Construir consulta dinámica con placeholders
	query := `
		SELECT idDoc, titulo, autor, categoria, formato, nombreArchivo
		FROM documentos 
		WHERE idDoc IN (` + placeholders(len(ides)) + `)
	`

	// Convertir IDs a interface{} para pasarlos como parámetros
	args := make([]interface{}, len(ides))
	for i, id := range ides {
		args[i] = id
	}

	// Ejecutar consulta
	rows, err := models.MariaDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf(" error al ejecutar la consulta en MariaDB: %w", err)
	}
	defer rows.Close()

	// Acumular resultados
	var recursos []*model.RecursoMuestra
	for rows.Next() {
		var recurso model.RecursoMuestra
		if err := rows.Scan(&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato, &recurso.Archivo); err != nil {
			return nil, fmt.Errorf("error  al leer los resultados de MariaDB: %w", err)
		}
		recursos = append(recursos, &recurso)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar  sobre los resultados de MariaDB: %w", err)
	}

	return recursos, nil
}

// ! tomarMuestraAleatoriaMariadb
func TomarMuestraAleatoriaMariadb(parametros model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	// Consulta en MariaDB
	rows, err := models.MariaDB.Query(`
		SELECT idDoc, titulo, autor, categoria, formato, nombreArchivo
		FROM documentos
		WHERE
			(titulo = ? OR ? IS NULL) AND
			(autor = ? OR ? IS NULL) AND
			(categoria = ? OR ? IS NULL) AND
			(formato = ? OR ? IS NULL) ORDER BY RAND()
			LIMIT ?;
	`, parametros.Titulo, parametros.Titulo, parametros.Autor, parametros.Autor, parametros.Categoria, parametros.Categoria, parametros.Formato, parametros.Formato, parametros.Cantidad)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en MariaDB: %w", err)
	}
	defer rows.Close()

	// Acumular resultados de MariaDB
	var recursos []*model.RecursoMuestra

	for rows.Next() {
		var recurso model.RecursoMuestra
		if err := rows.Scan(&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato, &recurso.Archivo); err != nil {
			return nil, fmt.Errorf("error al leer los resultados de MariaDB: %w", err)
		}
		recursos = append(recursos, &recurso)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados de MariaDB: %w", err)
	}
	return recursos, nil
}

// ! tomarMuestraUltimaMariadb
func TomarMuestraUltimaMariadb(parametros model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	// Consulta en MariaDB
	rows, err := models.MariaDB.Query(`
		SELECT idDoc, titulo, autor, categoria, formato, nombreArchivo
		FROM documentos
		WHERE
			(titulo = ? OR ? IS NULL) AND
			(autor = ? OR ? IS NULL) AND
			(categoria = ? OR ? IS NULL) AND
			(formato = ? OR ? IS NULL) ORDER BY idDoc DESC
			LIMIT 3;
	`, parametros.Titulo, parametros.Titulo, parametros.Autor, parametros.Autor, parametros.Categoria, parametros.Categoria, parametros.Formato, parametros.Formato)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en MariaDB: %w", err)
	}
	defer rows.Close()

	// Acumular resultados de MariaDB
	var recursos []*model.RecursoMuestra

	for rows.Next() {
		var recurso model.RecursoMuestra
		if err := rows.Scan(&recurso.ID, &recurso.Titulo, &recurso.Autor, &recurso.Categoria, &recurso.Formato, &recurso.Archivo); err != nil {
			return nil, fmt.Errorf("error al leer los resultados de MariaDB: %w", err)
		}
		recursos = append(recursos, &recurso)

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados de MariaDB: %w", err)
	}
	return recursos, nil
}

// ! tomarMuestraCouchbase
func TomarMuestraCouchbase(ides []int, parametros *model.ParametrosBusqueda) (map[string]map[string]interface{}, error) {
	// Verificar si la lista de IDs está vacía
	if len(ides) == 0 {
		return nil, fmt.Errorf("la lista de IDs está vacía")
	}

	// Construir los docIDs
	docIDes := make([]string, len(ides))
	for i, id := range ides {
		docIDes[i] = TYPE_DOC + strconv.Itoa(id)
	}

	// Construir la consulta base
	query := "SELECT META().id AS id, etiquetas, calificacion FROM `EduAndTime`.`documentos`.`documento` WHERE META().id IN $1"
	positionalParams := []interface{}{docIDes}

	// Agregar condición de etiquetas si es necesario

	if parametros != nil && parametros.Etiquetas != nil && len(parametros.Etiquetas) > 0 {
		query += " AND ANY etiqueta IN etiquetas SATISFIES etiqueta IN $2 END"
		positionalParams = append(positionalParams, parametros.Etiquetas)
	}

	// Ejecutar la consulta
	resultados, err := models.CouchbaseCluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: positionalParams,
	})
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}

	// Procesar los resultados de Couchbase
	documentoMap := make(map[string]map[string]interface{})
	for resultados.Next() {
		var resultado map[string]interface{}
		if err := resultados.Row(&resultado); err != nil {
			return nil, fmt.Errorf("error al leer los resultados de Couchbase: %w", err)
		}
		// Usar el ID del documento como clave
		if docID, ok := resultado["id"].(string); ok {
			documentoMap[docID] = resultado
		} else {
			return nil, fmt.Errorf("no se encontró el ID del documento en el resultado")
		}
	}
	if err := resultados.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los resultados de Couchbase: %w", err)
	}

	return documentoMap, nil
}
