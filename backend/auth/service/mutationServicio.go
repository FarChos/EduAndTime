package service

import (
	"ModuloAuth/db"
	"ModuloAuth/graph/model"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
)

// ! CrearEnMariadb crea el registro del ususario de mariadb.
func CrearEnMariadb(usuario model.UsuarioInput, contrasenaHash []byte) (int64, error) {
	// Obtener fecha actual en formato YYYY-MM-DD
	diaCreacion := time.Now().Format("2006-01-02")

	// Validar nombre de imagen
	nombreImagen := ""
	if usuario.Imagen != nil {
		nombreImagen = usuario.Imagen.Filename
	}

	// Ejecutar consulta para insertar el usuario
	result, err := db.MariaDB.Exec("INSERT INTO usuarios (nombre, correo, contrasena, nombreImagen, fechaOrigen) VALUES (?, ?, ?, ?, ?)",
		usuario.Nombre, usuario.Correo, contrasenaHash, nombreImagen, diaCreacion)
	if err != nil {
		return 0, fmt.Errorf("error al insertar el usuario en MariaDB: %v", err)
	}

	// Obtener ID del último registro insertado
	idDoc, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el ID del recurso insertado: %w", err)
	}

	return idDoc, nil
}

// ! Eliminar el registro del ususario de mariadb.
func EliminarEnMariadb(idUsuario int64) error {
	result, err := db.MariaDB.Exec("DELETE FROM usuarios WHERE idUsuario = ?", idUsuario)
	if err != nil {
		return fmt.Errorf("error al eliminar el usuario con ID %d: %v", idUsuario, err)
	}

	// Verificar si se afectó alguna fila
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error al verificar filas afectadas: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se encontró un usuario con ID %d para eliminar", idUsuario)
	}

	return nil
}

// ! CrearEnCouchbase crea un documento inicial para un usuario en Couchbase.
func CrearEnCouchbase(idUsuario int64) (bool, error) {
	// Crear el documento inicial para Couchbase
	doc := map[string]interface{}{
		"docFavoritos":   []int{},           // Lista vacía de favoritos
		"docCalificados": map[int64]int32{}, // Mapa vacío de calificados con int64 como clave y int como valor
		"docOriginados":  []int{},           // Lista vacía de documentos originados
	}

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Acceder al bucket y verificar disponibilidad
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return false, fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %v", err)
	}

	// Seleccionar la colección
	collection := bucket.Scope("usuarios").Collection("usuario")

	// Intentar insertar el documento
	_, err = collection.Insert(docID, doc, nil)
	if err != nil {
		return false, fmt.Errorf("error al insertar el usuario en Couchbase: %v", err)
	}

	return true, nil
}

// ! EliminarDeCouchbase eliminara el documento relacionado al usuario
func EliminarDeCouchbase(idUsuario int64) error {
	// Formatear la clave del documento que se va a eliminar
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Acceder al bucket y colección de Couchbase
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil) // Asegura que el bucket esté listo
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket en Couchbase: %v", err)
	}

	collection := bucket.Scope("usuarios").Collection("usuario")

	// Eliminar el documento en Couchbase
	_, err = collection.Remove(docID, nil)
	if err != nil {
		return fmt.Errorf("error al eliminar el documento en Couchbase: %v", err)
	}

	return nil
}

// ! ActualizarEnMariaDB actualiza el registro del usuario en mariadb.
func ActualizarEnMariaDB(id int64, usuario model.UsuarioInput, contrasenaHash []byte, nombreImagen string) error {
	// Verificar si el nombre es nil o vacío
	if usuario.Nombre == nil || *usuario.Nombre == "" {
		return fmt.Errorf("el nombre del usuario no puede estar vacío")
	}

	// Ejecutar la consulta
	_, err := db.MariaDB.Exec(
		"UPDATE usuarios SET nombre = ?, contrasena = ?, nombreImagen = ? WHERE idUsuario = ?",
		usuario.Nombre, contrasenaHash, nombreImagen, id,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar el usuario con ID %d: %v", id, err)
	}

	return nil // Retorno exitoso
}

// ! GuardarImagen guarda la imagen de perfil del usuario en la ruta /imagenesPerfil del docker
func GuardarImagen(usuario model.UsuarioInput, directorio string) (string, error) {
	// Verificar si el campo Imagen es nil
	if usuario.Imagen == nil {
		return "", fmt.Errorf("no se proporcionó una imagen")
	}

	// Obtener el archivo de la imagen (campo File de *graphql.Upload)
	archivo := usuario.Imagen.File
	if archivo == nil {
		return "", fmt.Errorf("el archivo de la imagen es nil")
	}

	// Definir el directorio donde guardar las imágenes
	// Ruta absoluta en el contenedor Docker

	// Verificar si el directorio existe, y si no, crear el directorio
	if _, err := os.Stat(directorio); os.IsNotExist(err) {
		err := os.MkdirAll(directorio, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("error al crear el directorio de imágenes: %v", err)
		}
	}

	// Crear un nombre único para la imagen (por ejemplo, usando la fecha y hora actual)
	nombreImagen := fmt.Sprintf("%d_%s", time.Now().Unix(), usuario.Imagen.Filename)

	// Establecer la ruta completa del archivo
	rutaArchivo := filepath.Join(directorio, nombreImagen)

	// Crear un archivo en el directorio con el nombre generado
	destino, err := os.Create(rutaArchivo)
	if err != nil {
		return "", fmt.Errorf("error al crear el archivo de imagen: %v", err)
	}
	defer destino.Close() // Cerrar el archivo destino cuando la función termine

	// Copiar el contenido del archivo de la imagen al nuevo archivo
	_, err = io.Copy(destino, archivo)
	if err != nil {
		return "", fmt.Errorf("error al guardar la imagen: %v", err)
	}

	// Devolver solo el nombre de la imagen guardada
	return nombreImagen, nil
}

// ! GuardarNuevoFavorito guarda el id de un recurso que el usuario marco como favorito en couchbase
func GuardarNuevoFavorito(idUsuario int64, idFavorito *int64) error {
	favoritoId := *idFavorito

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Acceder al bucket y la colección de Couchbase
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket en couchbase: %v", err)
	}
	collection := bucket.Scope("usuarios").Collection("usuario")

	// Utilizar MutateIn para agregar el nuevo favorito a `docFavoritos`
	mutateOps := []gocb.MutateInSpec{
		gocb.ArrayAddUniqueSpec(
			"docFavoritos", // Ruta del campo
			favoritoId,     // Valor único a agregar
			&gocb.ArrayAddUniqueSpecOptions{},
		),
	}

	// Ejecutar la operación MutateIn
	_, err = collection.MutateIn(docID, mutateOps, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrPathExists) {
			return fmt.Errorf("el favorito %d ya está en la lista", favoritoId)
		}
		return fmt.Errorf("error al agregar el nuevo favorito: %v", err)
	}

	fmt.Printf("Nuevo favorito %d agregado correctamente a docFavoritos\n", favoritoId)
	return nil
}

// ! GuardarNuevoAMisRecursos guarda el id de un recurso que el usuario subio al sistema en couchbase
func GuardarNuevoAMisRecursos(idUsuario int64, IdeMiRecurso *int64) error {
	// Validar el parámetro de entrada
	if IdeMiRecurso == nil {
		return fmt.Errorf("el ID del recurso no puede ser nulo")
	}
	miRecursoId := *IdeMiRecurso

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Acceder al bucket y la colección de Couchbase
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket: %v", err)
	}
	collection := bucket.Scope("usuarios").Collection("usuario")

	// Usar MutateIn para agregar el nuevo recurso a la lista docOriginados
	mutateOps := []gocb.MutateInSpec{
		gocb.ArrayAddUniqueSpec("docOriginados", miRecursoId, &gocb.ArrayAddUniqueSpecOptions{}),
	}

	_, err = collection.MutateIn(docID, mutateOps, nil)
	if err != nil {
		if errors.Is(err, gocb.ErrPathExists) {
			return fmt.Errorf("el recurso %d ya está en la lista", miRecursoId)
		}
		return fmt.Errorf("error al agregar el recurso al documento: %v", err)
	}

	fmt.Printf("Nuevo recurso %d agregado correctamente a docOriginados\n", miRecursoId)
	return nil
}

// ! GuardarRecursoCalificado guarda la calificación y el id que el usuario le dio a un recurso en couchbase
func GuardarRecursoCalificado(id int64, recursoCalificado *model.RecursoCalificadoInput) error {
	if recursoCalificado == nil {
		return fmt.Errorf("recursoCalificado no puede ser nil")
	}

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(id, 10)

	// Acceder al bucket y la colección de Couchbase
	bucket := db.CouchbaseCluster.Bucket("EduAndTime")
	err := bucket.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del bucket: %v", err)
	}
	collection := bucket.Scope("usuarios").Collection("usuario")

	// Usar MutateIn para actualizar solo `docCalificados`
	mutateOps := []gocb.MutateInSpec{
		gocb.UpsertSpec(
			"docCalificados."+strconv.FormatInt(recursoCalificado.IDCalificado, 10), // Ruta de la clave en el campo
			recursoCalificado.Calificacion,                                          // Nuevo valor
			nil,                                                                     // Opciones por defecto
		),
	}

	_, err = collection.MutateIn(docID, mutateOps, nil)
	if err != nil {
		return fmt.Errorf("error al actualizar docCalificados: %v", err)
	}

	fmt.Printf("Recurso calificado con ID %d agregado/actualizado correctamente con calificación %d\n", recursoCalificado.IDCalificado, recursoCalificado.Calificacion)
	return nil
}

// ! EliminarFavorito
func EliminarFavorito(idUsuario int64, idFavorito *int64) error {
	// Asegurarse de que el idFavorito no sea nil
	if idFavorito == nil {
		return fmt.Errorf("el idFavorito no puede ser nulo")
	}

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Obtener el cluster
	cluster := db.CouchbaseCluster
	err := cluster.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del cluster : %v", err)
	}

	// Consulta N1QL para eliminar el valor de la lista usando la clave del documento
	stmt := "UPDATE `EduAndTime` SET docFavoritos = ARRAY_REMOVE(docFavoritos, ?) WHERE META().id = ?"

	// Ejecutar la consulta N1QL usando cluster.Query, pasando los parámetros a través de PositionalParameters
	queryOptions := &gocb.QueryOptions{
		PositionalParameters: []interface{}{*idFavorito, docID},
	}

	// Ejecutar la consulta
	rows, err := cluster.Query(stmt, queryOptions)
	if err != nil {
		return fmt.Errorf("error al eliminar el favorito usando N1QL : %w", err)
	}

	// Esperar a que se procesen los resultados (si es necesario)
	defer rows.Close()

	return nil
}

// !
func EliminarDeMisRecursos(idUsuario int64, idMiRecurso *int64) error {
	// Asegurarse de que el idUsuario no sea nil
	if idMiRecurso == nil {
		return fmt.Errorf("el idFavorito no puede ser nil")
	}

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Obtener el cluster
	cluster := db.CouchbaseCluster
	err := cluster.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del cluster: %v", err)
	}

	// Consulta N1QL para eliminar el valor de la lista usando la clave del documento
	stmt := "UPDATE `EduAndTime` SET docOriginados = ARRAY_REMOVE(docOriginados, ?) WHERE META().id = ?"

	// Ejecutar la consulta N1QL usando cluster.Query, pasando los parámetros a través de PositionalParameters
	queryOptions := &gocb.QueryOptions{
		PositionalParameters: []interface{}{*idMiRecurso, docID},
	}

	// Ejecutar la consulta
	rows, err := cluster.Query(stmt, queryOptions)
	if err != nil {
		return fmt.Errorf("error al eliminar el favorito usando N1QL: %w", err)
	}

	// Esperar a que se procesen los resultados (si es necesario)
	defer rows.Close()

	return nil
}
func EliminarRecursoCalificado(idUsuario int64, recursoCalificado *model.RecursoCalificadoInput) error {
	// Asegurarse de que el recurso no sea nil
	if recursoCalificado == nil {
		return fmt.Errorf("el recurso no puede ser nil")
	}

	// Validar que recursoCalificado.ID no esté vacío
	if recursoCalificado.IDCalificado == 0 {
		return fmt.Errorf("el ID del recurso calificado no puede ser 0")
	}

	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Obtener el cluster
	cluster := db.CouchbaseCluster
	err := cluster.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return fmt.Errorf("error al esperar disponibilidad del cluster: %v", err)
	}

	// Consulta N1QL para eliminar el valor de la lista usando la clave del documento
	stmt := "UPDATE `EduAndTime` SET docCalificados = ARRAY_REMOVE(docCalificados, ?) WHERE META().id = ?"

	// Ejecutar la consulta N1QL usando cluster.Query, pasando los parámetros a través de PositionalParameters
	queryOptions := &gocb.QueryOptions{
		PositionalParameters: []interface{}{recursoCalificado.IDCalificado, docID},
	}

	// Ejecutar la consulta
	rows, err := cluster.Query(stmt, queryOptions)
	if err != nil {
		return fmt.Errorf("error al eliminar el recurso calificado usando N1QL: %w", err)
	}

	// Esperar a que se procesen los resultados (si es necesario)
	defer rows.Close()

	return nil
}
