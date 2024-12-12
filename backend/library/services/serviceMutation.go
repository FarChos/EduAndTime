package services

import (
	"fmt"
	"io"
	"libreria/graph/model"
	"libreria/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
)

// ! Actualiza el numero de descargas en couchbase
func ActualizarNoDescargas(noDescargas int, docID int) error {
	queryUpdate := "UPDATE `EduAndTime`.`documentos`.`documento` SET noDescargas = $1 WHERE META().id = $2"
	_, err := models.CouchbaseCluster.Query(queryUpdate, &gocb.QueryOptions{
		PositionalParameters: []interface{}{noDescargas, docID},
	})
	if err != nil {
		return fmt.Errorf("error al ejecutar la consulta de actualización: %w", err)
	}

	return nil
}
func ActualizarCalificacion(resultadoFinal []float64, id int) error {
	docID := "documento::" + strconv.Itoa(id)
	// Actualizar en Couchbase
	newCalificacion := resultadoFinal[1]
	newNoCalificaciones := int(resultadoFinal[2])
	queryUpdate := "UPDATE `EduAndTime`.`documentos`.`documento` SET calificacion = $1, noCalificaciones = $2 WHERE META().id = $3"
	_, err := models.CouchbaseCluster.Query(queryUpdate, &gocb.QueryOptions{
		PositionalParameters: []interface{}{newCalificacion, newNoCalificaciones, docID},
	})
	if err != nil {
		return fmt.Errorf("error al ejecutar la consulta de actualización: %w", err)
	}

	return nil
}

// ! crear el recurso en la base de datos mariadb
func CrearRecursoEnMariadb(recurso model.RecursoInput, nombreArchivo string) (int64, error) {
	// Validar campos requeridos
	if recurso.Titulo == "" || recurso.Autor == "" || recurso.Categoria == "" || recurso.IDUsuario == 0 || recurso.Formato == "" {
		return 0, fmt.Errorf("los campos Titulo, Autor, Categoria, IDUsuario y Formato son obligatorios")
	}

	// Generar fecha de origen
	fechaOrigen := time.Now().Format("2006-01-02")

	// Ejecutar consulta de inserción
	result, err := models.MariaDB.Exec(
		`INSERT INTO documentos 
		(titulo, autor, categoria, idUsuario, formato, descripcion, fechaOrigen, nombreArchivo) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		recurso.Titulo, recurso.Autor, recurso.Categoria, recurso.IDUsuario, recurso.Formato, recurso.Descripcion, fechaOrigen, nombreArchivo,
	)
	if err != nil {
		return 0, fmt.Errorf("error al insertar el recurso en MariaDB: %w", err)
	}

	// Obtener el ID del último registro insertado
	idDoc, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el ID del recurso insertado: %w", err)
	}

	return idDoc, nil
}
func CrearRecursoEnCouchbase(idDoc int64, etiquetas []string) error {
	// Validar entradas
	if idDoc <= 0 {
		return fmt.Errorf("el ID del documento debe ser mayor que 0, recibido: %d", idDoc)
	}

	if etiquetas == nil {
		etiquetas = []string{} // Asegurar que no sea nil
	}

	// Crear el documento
	doc := map[string]interface{}{
		"noCalificaciones": 0,
		"calificacion":     0.0,
		"noDescargas":      0,
		"etiquetas":        etiquetas,
	}

	// Crear el ID del documento
	docID := "documento::" + strconv.FormatInt(idDoc, 10)

	// Inicializar el bucket
	bucket := models.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket 'EduAndTime': %w", err)
	}

	// Obtener la colección
	collection := bucket.Scope("documentos").Collection("documento")

	// Insertar el documento
	_, err = collection.Insert(docID, doc, nil)
	if err != nil {
		// Considerar revertir en MariaDB si se usa en combinación
		return fmt.Errorf("error al insertar el documento en Couchbase con ID '%s': %w", docID, err)
	}

	return nil
}

func GuardarRecurso(recurso model.RecursoInput, directorio string) (string, error) {
	// Validar que el recurso tenga un nombre de archivo
	if recurso.Recurso.Filename == "" {
		return "", fmt.Errorf("no se proporcionó un nombre de recurso")
	}

	// Obtener el archivo desde el recurso
	archivo := recurso.Recurso.File
	if archivo == nil {
		return "", fmt.Errorf("el archivo del recurso es nil")
	}

	// Validar o crear el directorio si no existe
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err := os.MkdirAll(directorio, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error al crear el directorio '%s': %v", directorio, err)
		}
	}

	// Generar un nombre único para el recurso
	nombreRecurso := recurso.Recurso.Filename
	rutaArchivo := filepath.Join(directorio, nombreRecurso)

	// Prevenir colisión de nombres (opcional)
	if _, err := os.Stat(rutaArchivo); err == nil {
		nombreRecurso = fmt.Sprintf("%d_%s", time.Now().Unix(), recurso.Recurso.Filename)
		rutaArchivo = filepath.Join(directorio, nombreRecurso)
	}

	// Crear un archivo en el directorio
	destino, err := os.Create(rutaArchivo)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo '%s': %v", rutaArchivo, err)
	}
	defer func() {
		if cerr := destino.Close(); cerr != nil {
			fmt.Printf("error al cerrar el archivo destino: %v\n", cerr)
		}
	}()

	// Copiar el contenido del archivo
	_, err = io.Copy(destino, archivo)
	if err != nil {
		return "", fmt.Errorf("error al guardar el archivo '%s': %v", rutaArchivo, err)
	}

	// Devolver el nombre del archivo guardado
	return nombreRecurso, nil
}

func EliminarRecursoEnMariadb(idDoc int64) error {
	result, err := models.MariaDB.Exec("DELETE FROM documentos WHERE idDoc = ?", idDoc)
	if err != nil {
		return fmt.Errorf("error al eliminar el documento con ID %d: %v", idDoc, err)
	}

	// Verificar si se afectó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró un documento con ID %d para eliminar", idDoc)
	}

	return nil
}

func EliminarDeCouchbase(idUsuario int64) error {
	// Formatear la clave del documento que se va a eliminar
	docID := "documento::" + strconv.FormatInt(idUsuario, 10)

	// Acceder al bucket y colección de Couchbase
	bucket := models.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil) // Asegura que el bucket esté listo
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %v", err)
	}

	collection := bucket.Scope("documentos").Collection("documento")

	// Eliminar el documento en Couchbase
	_, err = collection.Remove(docID, nil)
	if err != nil {
		return fmt.Errorf("error al eliminar el documento en Couchbase: %v", err)
	}

	return nil
}
