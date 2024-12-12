package controller

import (
	"ModuloAuth/graph/model"
	"ModuloAuth/service"
)

// !
func AutentificarUsuario(input model.UsuarioInput) (model.TokenAuth, error) {
	tokenAuth := model.TokenAuth{Exito: false}
	err := service.AutentificarUsuario(input.Contrasena, input.Correo)
	if err != nil {
		return tokenAuth, err
	}
	usuario, err := service.TomarDatosUsuarioDeMariadb(input.Correo)
	if err != nil {
		return tokenAuth, err
	}
	token, err := service.GenerarToken(usuario.Nombre)
	if err != nil {
		return tokenAuth, err
	}
	tokenAuth.Usuario = &usuario
	tokenAuth.Token = &token
	tokenAuth.Exito = true
	return tokenAuth, nil
}

// !
func TomarUsuarioXRecursos(idUsuario int64) (model.UsuarioXRecursos, error) {
	var DatosRecursoXUsuario model.UsuarioXRecursos
	rowsRecursoXUsuario, err := service.TomarUsuarioXRecursosDeCouchbase(idUsuario)
	if err != nil {
		return DatosRecursoXUsuario, err
	}
	ClaveValor, err := service.DevolverResultadoClaveValor(rowsRecursoXUsuario)
	if err != nil {
		return DatosRecursoXUsuario, err
	}
	DatosRecursoXUsuario, err = service.DevolverUsuarioXRecursos(ClaveValor)
	return DatosRecursoXUsuario, nil
}
