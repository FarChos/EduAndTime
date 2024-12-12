package controller

import (
	"ModuloAuth/graph/model"
	"ModuloAuth/service"
	"fmt"
)

const RUTA_IMAGENES_PERFIL string = "/imagenesPerfil"

// ! Crea al usuario en las bases de datos
func CrearUsuario(usuario model.UsuarioInput) (model.Resultado, error) {
	resultado := model.Resultado{
		Exito: false,
	}

	// Generar el hash de la contraseña
	contrasenaHash, err := service.HashContrasena(usuario.Contrasena)
	if err != nil {
		return resultado, fmt.Errorf("error al hashear la contraseña: %v", err)
	}

	// Crear el usuario en MariaDB
	idUsuarioCreado, err := service.CrearEnMariadb(usuario, contrasenaHash)
	if err != nil {
		return resultado, fmt.Errorf("error al crear usuario en MariaDB: %v", err)
	}

	// Crear el documento en Couchbase
	esCreadoEnCouchbase, err := service.CrearEnCouchbase(idUsuarioCreado)
	if !esCreadoEnCouchbase {
		// Rollback en MariaDB si falla Couchbase
		err := service.EliminarEnMariadb(idUsuarioCreado)
		if err != nil {
			return resultado, fmt.Errorf("error al eliminar al usuario en mariadb: %v", err)
		}
		return resultado, fmt.Errorf("error al crear el documento en Couchbase: %v", err)
	}

	// Actualizar el resultado a éxito
	resultado.Exito = true
	return resultado, nil
}

// ! Actualiza los algunos datos del usuario
func ActualizarUsuario(usuario model.UsuarioInput) (model.Resultado, error) {
	// Inicializar resultado
	resultado := model.Resultado{
		Exito: false,
	}

	// Verificar si la contraseña está vacía (si no se quiere actualizar)
	var contrasenaHash []byte
	if usuario.Contrasena != "" {
		// Generar el hash de la contraseña
		var err error
		contrasenaHash, err = service.HashContrasena(usuario.Contrasena)
		if err != nil {
			return resultado, fmt.Errorf("error al hashear la contraseña para el usuario con correo %s: %v", usuario.Correo, err)
		}
	}

	// Buscar el ID del usuario usando el correo
	idUsuario, err := service.BuscarIdConCorreo(usuario.Correo)
	if err != nil {
		return resultado, fmt.Errorf("error al buscar el usuario con correo %s en la base de datos: %v", usuario.Correo, err)
	}
	nombreImagen, err := service.GuardarImagen(usuario, RUTA_IMAGENES_PERFIL)
	if err != nil {
		return resultado, fmt.Errorf("error al guardar la imagen de perfil del ususario %s en la ruta /imagenesPerfil: %v", usuario.Correo, err)
	}
	// Actualizar usuario en MariaDB
	err = service.ActualizarEnMariaDB(idUsuario, usuario, contrasenaHash, nombreImagen)
	if err != nil {
		return resultado, fmt.Errorf("error al actualizar el usuario con ID %d en MariaDB: %v", idUsuario, err)
	}

	// Retornar resultado exitoso
	resultado.Exito = true
	return resultado, nil
}

// ! Eliminación completa de los datos de usuario en ambas bases de datos
func EliminarUsuario(usuario model.UsuarioInput) (model.Resultado, error) {
	resultado := model.Resultado{
		Exito: false,
	}

	// Autentificación del usuario
	err := service.AutentificarUsuario(usuario.Contrasena, usuario.Correo)
	if err != nil {
		return resultado, fmt.Errorf("error de autenticación: %v", err)
	}

	// Buscar el ID del usuario en la base de datos
	idUsuario, err := service.BuscarIdConCorreo(usuario.Correo)
	if err != nil {
		return resultado, fmt.Errorf("error al buscar el usuario en la base de datos: %w", err)
	}

	// Eliminar el usuario de MariaDB
	err = service.EliminarEnMariadb(idUsuario)
	if err != nil {
		return resultado, fmt.Errorf("error al eliminar el usuario en MariaDB: %w", err)
	}

	// Eliminar el usuario de Couchbase
	err = service.EliminarDeCouchbase(idUsuario)
	if err != nil {
		return resultado, fmt.Errorf("error al eliminar el usuario en Couchbase: %w", err)
	}

	// Si ambas operaciones fueron exitosas, actualizar el resultado
	resultado.Exito = true

	return resultado, nil
}

// ! ActualizarUsuarioXRecursos actualiza los datos del usuario que estan relacionados a los recursos en couchbase
func ActualizarUsuarioXRecursos(id int64, input model.UsuarioXRecursoInput) (model.Resultado, error) {
	mensaje := "operación realizada con éxito"
	resultado := model.Resultado{
		Exito: false,
	}

	// Verificación explícita de que al menos un campo tiene valor
	if input.IdeFavorito == nil && input.IdeMiRecurso == nil && input.RecursoCalificado == nil {
		return resultado, fmt.Errorf("Todos los datos están vacíos")
	}

	// Procesar el campo IdeFavorito
	if input.IdeFavorito != nil {
		err := service.GuardarNuevoFavorito(id, input.IdeFavorito)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo favorito : %w", err)
		}
	}

	// Procesar el campo IdeMiRecurso
	if input.IdeMiRecurso != nil {
		err := service.GuardarNuevoAMisRecursos(id, input.IdeMiRecurso)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo recurso propio: %w", err)
		}
	}

	// Procesar el campo RecursoCalificado
	if input.RecursoCalificado != nil {
		err := service.GuardarRecursoCalificado(id, input.RecursoCalificado)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo documento calificado: %w", err)
		}
	}

	// Si ha ido bien, se marca como exitoso
	resultado.Exito = true
	resultado.Mensaje = &mensaje
	return resultado, nil
}

// ! EliminarUsuarioXRecurso elimina los datos del usuario que estan relacionados a los recursos en couchbase
func EliminarUsuarioXRecurso(id int64, input model.UsuarioXRecursoInput) (model.Resultado, error) {
	mensaje := "operación realizada con éxito"
	resultado := model.Resultado{
		Exito: false,
	}

	// Verificación explícita de que al menos un campo tiene valor
	if input.IdeFavorito == nil && input.IdeMiRecurso == nil && input.RecursoCalificado == nil {
		return resultado, fmt.Errorf("Todos los datos están vacíos")
	}

	// Procesar el campo IdeFavorito
	if input.IdeFavorito != nil {
		err := service.EliminarFavorito(id, input.IdeFavorito)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo favorito : %w", err)
		}
	}

	// Procesar el campo IdeMiRecurso
	if input.IdeMiRecurso != nil {
		err := service.EliminarDeMisRecursos(id, input.IdeMiRecurso)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo recurso propio: %w", err)
		}
	}

	// Procesar el campo RecursoCalificado
	if input.RecursoCalificado != nil {
		err := service.EliminarRecursoCalificado(id, input.RecursoCalificado)
		if err != nil {
			return resultado, fmt.Errorf("error al guardar nuevo documento calificado: %w", err)
		}
	}

	// Si ha ido bien, se marca como exitoso
	resultado.Exito = true
	resultado.Mensaje = &mensaje
	return resultado, nil

}
