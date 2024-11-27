package services

import (
	"fmt"
	"libreria/graph/model"
	"libreria/models"
	"strconv"

	"github.com/couchbase/gocb/v2"
)

func convertirEtiquetas(etiquetas []*string) []string {
	var resultado []string
	for _, etiqueta := range etiquetas {
		if etiqueta != nil { // Asegurarte de que no sea nil
			resultado = append(resultado, *etiqueta)
		}
	}
	return resultado
}

func devolverResultadoClaveValor(rows *gocb.QueryResult) (map[string]interface{}, error) {
	var resultado map[string]interface{}

	// Iteramos sobre las filas con `rows.Rows()`
	if rows.Next() {
		// Recuperamos la fila actual
		err := rows.Row(&resultado)
		if err != nil {
			return nil, fmt.Errorf("error al leer la fila: %w", err)
		}
		// Como esperamos un solo resultado, podemos retornar aquí
		return resultado, nil
	}

	// Si no se encuentran filas
	if rows.Err() != nil {
		return nil, fmt.Errorf("error en la iteración de las filas: %w", rows.Err())
	}
	return nil, fmt.Errorf("no se encontró el documento")
}

func manejoEtiquetas(etiquetasRaw []interface{}) ([]*string, error) {
	if len(etiquetasRaw) == 0 {
		// Si no hay etiquetas, se devuelve nil, sin error
		return nil, nil
	}

	// Si hay etiquetas, se procesan
	etiquetas := make([]*string, len(etiquetasRaw))
	for i, etiqueta := range etiquetasRaw {
		etiquetaStr, ok := etiqueta.(string)
		if !ok {
			return nil, fmt.Errorf("una de las etiquetas no es un string")
		}
		etiquetas[i] = &etiquetaStr
	}
	return etiquetas, nil
}

func devolverDatosCouchbase(recurso model.RecursoMuestra) (model.RecursoMuestra, error) {
	docID := "documento::" + strconv.Itoa(recurso.ID)

	query := "SELECT etiquetas, calificacion FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
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
	return recurso, nil
}
