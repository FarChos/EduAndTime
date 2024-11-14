package graph

import (
	"ModuloAuth/controller"
	"ModuloAuth/graph/model"
	"context"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func (r *queryResolver) AutentificarUsuario(ctx context.Context, correo string, contrasena string) (model.TokenAuth, error) {
	tokenAuth, err := controller.AutentificarAlUsuario(correo, contrasena)
	if err != nil {
		tokenAuth = model.TokenAuth{Exito: false}
		return tokenAuth, err
	}
	return tokenAuth, nil
}

func (r mutationResolver) CrearUsuario(ctx context.Context, nombre string, correo string, imgPerf *string, contrasena string) (bool, error) {

	usuario := &model.Usuario{
		Nombre:  nombre,
		Correo:  correo,
		ImgPerf: imgPerf,
	}
	esCreado, err := controller.CrearAlUsuario(usuario, contrasena)
	if err != nil {
		return esCreado, err
	}

	return esCreado, nil
}

func (r mutationResolver) ActualizarUsuario(ctx context.Context, correo string, nombre string, contrasena string) (bool, error) {

	usuario := &model.Usuario{
		Nombre: nombre,
		Correo: correo,
	}
	esActualizado, err := controller.ActualizarAlUsuario(usuario, contrasena)
	if err != nil {
		return esActualizado, err
	}

	return esActualizado, nil
}
func (r mutationResolver) EliminarUsuario(ctx context.Context, correo string) (bool, error) {

	esEliminado, err := controller.EliminarAlUsuario(correo)
	if err != nil {
		return esEliminado, err
	}

	return esEliminado, nil
}
