package controller

import (
	"ModuloAuth/graph/model"
	"ModuloAuth/service"
)

func AutentificarAlUsuario(correo string, contrasena string) (model.TokenAuth, error) {
	tokenAuth, err := service.BuscarUsuario(correo, contrasena)
	if err != nil {
		return tokenAuth, err
	}
	return tokenAuth, nil
}
func CrearAlUsuario(usuario *model.Usuario, contrasena string) (bool, error) {
	esCreado, err := service.CrearUsuario(usuario, contrasena)
	if err != nil {
		return esCreado, err
	}
	return esCreado, nil
}
func ActualizarAlUsuario(usuario *model.Usuario, contrasena string) (bool, error) {
	esActualizado, err := service.ActualizarUsuario(usuario, contrasena)
	if err != nil {
		return esActualizado, err
	}
	return esActualizado, err
}

func EliminarAlUsuario(correo string) (bool, error) {
	esEliminado, err := service.EliminarUsuario(correo)
	if err != nil {
		return esEliminado, err
	}
	return esEliminado, err
}
