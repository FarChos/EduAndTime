package services

import (
	"errors"
	"fmt"
	"libreria/graph/model"
	"libreria/models"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
)

func ActualizarNumDescargas(id int) (bool, error) {
	// Construir el ID del documento
	docID := "documento::" + strconv.FormatInt(int64(id), 10)

	// Consulta para obtener noDescargas
	queryRead := "SELECT noDescargas FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(queryRead, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return false, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}

	// Procesar el resultado de la consulta
	resultado, err := devolverResultadoClaveValor(rows)
	if err != nil {
		return false, fmt.Errorf("error al procesar el resultado de la consulta: %w", err)
	}

	// Extraer y verificar el valor de noDescargas
	noDescargasRaw, ok := resultado["noDescargas"]
	if !ok {
		return false, fmt.Errorf("el campo 'noDescargas' no existe en el documento con ID %s", docID)
	}

	// Convertir a float64 y luego a int
	noDescargasFloat, ok := noDescargasRaw.(float64)
	if !ok {
		return false, fmt.Errorf("el campo 'noDescargas' no es del tipo esperado en el documento con ID %s", docID)
	}
	noDescargas := int(noDescargasFloat) + 1

	// Consulta para actualizar noDescargas
	queryUpdate := "UPDATE `EduAndTime`.`documentos`.`documento` SET noDescargas = $1 WHERE META().id = $2"
	_, err = models.CouchbaseCluster.Query(queryUpdate, &gocb.QueryOptions{
		PositionalParameters: []interface{}{noDescargas, docID},
	})
	if err != nil {
		return false, fmt.Errorf("error al ejecutar la consulta de actualización: %w", err)
	}

	return true, nil
}

func CalificarRecurso(id int, calificacion float32) (bool, error) {
	// Construir el ID del documento
	docID := "documento::" + strconv.Itoa(id)

	// Consulta para obtener noCalificaciones y calificacion
	queryRead := "SELECT noCalificaciones, calificacion FROM `EduAndTime`.`documentos`.`documento` WHERE META().id = $1"
	rows, err := models.CouchbaseCluster.Query(queryRead, &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	})
	if err != nil {
		return false, fmt.Errorf("error al ejecutar la consulta en Couchbase: %w", err)
	}

	// Procesar el resultado de la consulta
	resultado, err := devolverResultadoClaveValor(rows)
	if err != nil {
		return false, fmt.Errorf("error al procesar el resultado de la consulta: %w", err)
	}

	// Manejar noCalificaciones
	noCalificacionesRaw, ok := resultado["noCalificaciones"].(float64)
	if !ok {
		noCalificacionesRaw = 0 // Si no existe, asumimos 0 calificaciones
	}

	// Manejar calificacion
	calificacionActualRaw, ok := resultado["calificacion"].(float64)
	if !ok {
		calificacionActualRaw = 0 // Si no existe, asumimos calificación inicial 0
	}

	// Convertir valores para cálculos
	noCalificaciones := int(noCalificacionesRaw)
	calificacionActual := float32(calificacionActualRaw)

	// Calcular nueva calificación
	newCalificacion := (calificacionActual*float32(noCalificaciones) + calificacion) / float32(noCalificaciones+1)
	newNoCalificaciones := noCalificaciones + 1

	// Actualizar en Couchbase
	queryUpdate := "UPDATE `EduAndTime`.`documentos`.`documento` SET calificacion = $1, noCalificaciones = $2 WHERE META().id = $3"
	_, err = models.CouchbaseCluster.Query(queryUpdate, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newCalificacion, newNoCalificaciones, docID},
	})
	if err != nil {
		return false, fmt.Errorf("error al ejecutar la consulta de actualización: %w", err)
	}

	return true, nil
}

func CrearRecurso(recurso model.RecursoInput) (bool, error) {
	// Validar los datos de entrada si es necesario

	fechaOrigen := time.Now().Format("2006-01-02")
	result, err := models.MariaDB.Exec(
		"INSERT INTO documentos (titulo, autor, categoria, idUsuario, formato, descripcion, fechaOrigen, direcDoc) VALUES (?,?,?,?,?,?,?,?)",
		recurso.Titulo, recurso.Autor, recurso.Categoria, recurso.IDUsuario, recurso.Formato, recurso.Descripcion, fechaOrigen, recurso.DireccionRec,
	)
	if err != nil {
		return false, fmt.Errorf("error al insertar el recurso en MariaDB: %w", err)
	}

	idDoc, err := result.LastInsertId()
	if err != nil {
		return false, fmt.Errorf("error al obtener el ID del recurso insertado: %w", err)
	}

	doc := map[string]interface{}{
		"noCalificaciones": 0,
		"calificacion":     nil,
		"noDescargas":      0,
		"etiquetas":        convertirEtiquetas(recurso.Etiquetas),
	}

	docID := "documento::" + strconv.FormatInt(idDoc, 10)

	bucket := models.CouchbaseCluster.Bucket("EduAndTime")
	err = bucket.WaitUntilReady(10*time.Second, nil) // Aumentar tiempo si necesario
	if err != nil {
		return false, fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %w", err)
	}

	collection := bucket.Scope("documentos").Collection("documento")
	_, err = collection.Insert(docID, doc, nil)
	if err != nil {
		// Considerar revertir la inserción en MariaDB si es necesario
		return false, fmt.Errorf("error al insertar el usuario en Couchbase: %w", err)
	}

	return true, nil
}

func EliminarRecurso(id int) (bool, error) {
	// Eliminar el recurso en MariaDB
	result, err := models.MariaDB.Exec("DELETE FROM documentos WHERE idDoc = ?", id)
	if err != nil {
		return false, fmt.Errorf("error al eliminar el recurso en MariaDB: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error al verificar eliminación en MariaDB: %v", err)
	}
	if rowsAffected == 0 {
		return false, fmt.Errorf("el recurso con ID %d no existe en MariaDB", id)
	}

	// Formar el docID para Couchbase
	docID := "documento::" + strconv.FormatInt(int64(id), 10)

	// Acceder al bucket y colección de Couchbase
	bucket := models.CouchbaseCluster.Bucket("EduAndTime")
	err = bucket.WaitUntilReady(10*time.Second, nil) // Asegurarse de que el bucket esté listo
	if err != nil {
		return false, fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %w", err)
	}

	collection := bucket.Scope("documentos").Collection("documento")

	// Eliminar el documento en Couchbase
	_, err = collection.Remove(docID, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			// Si el documento no existe en Couchbase, es un error no crítico
			return false, fmt.Errorf("el documento con ID %s no existe en Couchbase", docID)
		}
		return false, fmt.Errorf("error al eliminar el documento con ID %s en Couchbase: %w", docID, err)
	}

	return true, nil
}
