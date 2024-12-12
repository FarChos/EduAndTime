package service

import (
	"ModuloAuth/db"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"ModuloAuth/graph/model"

	"github.com/couchbase/gocb/v2"
	"golang.org/x/crypto/bcrypt"
)

// ! AutentificarUsuario asegura que el usuario ingreso una contraseña correcta
func AutentificarUsuario(contrasena string, correo string) error {
	var contrasenaAlmacenada string

	// Ejecutar la consulta para obtener la contraseña almacenada
	err := db.MariaDB.QueryRow("SELECT contrasena FROM usuarios WHERE correo = ?", correo).Scan(&contrasenaAlmacenada)
	if err != nil {
		// Si no se encuentra el usuario, el error es generalmente 'sql.ErrNoRows'
		if err == sql.ErrNoRows {
			return fmt.Errorf("usuario no encontrado")
		}
		// En caso de otros errores (como problemas de conexión), devolvemos el error
		return fmt.Errorf("error al consultar la base de datos: %v", err)
	}

	// Convertir las contraseñas a byte slices para la comparación
	contrasenaAlmacenadaByte := []byte(contrasenaAlmacenada)
	contrasenaByte := []byte(contrasena)

	// Comparar las contraseñas
	err = bcrypt.CompareHashAndPassword(contrasenaAlmacenadaByte, contrasenaByte)
	if err != nil {
		// Si la comparación falla, la contraseña es incorrecta
		return fmt.Errorf("contraseña incorrecta")
	}

	// Si va bien, las contraseñas coinciden
	return nil
}

// ! Retorna el id del usuario de la base de datos maria db usando su correo como parametro de busqueda
func BuscarIdConCorreo(correo string) (int64, error) {
	if correo == "" {
		return 0, fmt.Errorf("el correo no puede estar vacío")
	}

	var id int64
	err := db.MariaDB.QueryRow("SELECT idUsuario FROM usuarios WHERE correo = ?", correo).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no se encontró un usuario con el correo: %s", correo)
		}
		return 0, fmt.Errorf("error al buscar el ID para el correo %s: %v", correo, err)
	}

	return id, nil
}

// ! TomarDatosUsuarioDeMariadb toma los datos basicos de ususario para su inicio de sesión
func TomarDatosUsuarioDeMariadb(correo string) (model.Usuario, error) {
	var usuario model.Usuario

	// Realizar la consulta
	err := db.MariaDB.QueryRow(
		"SELECT idUsuario, nombre, COALESCE(nombreImagen, '') AS nombreImagen FROM usuarios WHERE correo = ?",
		correo,
	).Scan(&usuario.ID, &usuario.Nombre, &usuario.NombreImagen)

	// Manejo de errores
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Retornar un error personalizado si el usuario no se encuentra
			return usuario, fmt.Errorf("usuario no encontrado para el correo: %s", correo)
		}
		// Retornar cualquier otro error
		return usuario, fmt.Errorf("error al consultar el usuario: %w", err)
	}

	// Asignar el correo después de asegurarse de que no hay errores
	usuario.Correo = correo
	return usuario, nil
}

// ! UsuarioXRecursos
func TomarUsuarioXRecursosDeCouchbase(idUsuario int64) (*gocb.QueryResult, error) {
	// Generar la clave del documento
	docID := "usuario::" + strconv.FormatInt(idUsuario, 10)

	// Obtener el cluster
	cluster := db.CouchbaseCluster
	err := cluster.WaitUntilReady(10*time.Second, nil)
	if err != nil {
		return nil, fmt.Errorf("error al esperar disponibilidad del cluster: %v", err)
	}

	// Consulta N1QL para obtener los campos específicos
	stmt := `
		SELECT docFavoritos, docCalificados, docOriginados
		FROM EduAndTime
		WHERE META().id = $1
	`

	// Ejecutar la consulta con el parámetro del ID del documento
	queryOptions := &gocb.QueryOptions{
		PositionalParameters: []interface{}{docID},
	}
	rows, err := cluster.Query(stmt, queryOptions)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar la consulta N1QL: %w", err)
	}

	// Devolver las filas para que el consumidor procese y cierre
	return rows, nil
}
