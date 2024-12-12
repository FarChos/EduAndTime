package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"ModuloAuth/controller"
	"ModuloAuth/graph/model"
	"context"
)

type Resolver struct{}

// ActualizarUsuarioXRecursos is the resolver for the actualizarUsuarioXRecursos field.
func (r *mutationResolver) ActualizarUsuarioXRecursos(ctx context.Context, id int64, input model.UsuarioXRecursoInput) (model.Resultado, error) {
	resultado, err := controller.ActualizarUsuarioXRecursos(id, input)
	if err != nil {
		return resultado, err
	}

	return resultado, nil
}

// EliminarUsuarioXRecurso is the resolver for the eliminarUsuarioXRecurso field.
func (r *mutationResolver) EliminarUsuarioXRecurso(ctx context.Context, id int64, input model.UsuarioXRecursoInput) (model.Resultado, error) {
	resultado, err := controller.EliminarUsuarioXRecurso(id, input)
	if err != nil {
		return resultado, err
	}

	return resultado, nil
}

// CrearUsuario is the resolver for the crearUsuario field.
func (r *mutationResolver) CrearUsuario(ctx context.Context, input model.UsuarioInput) (model.Resultado, error) {
	resultado, err := controller.CrearUsuario(input)
	if err != nil {
		return resultado, err
	}

	return resultado, nil
}

// ActualizarUsuario is the resolver for the actualizarUsuario field.
func (r *mutationResolver) ActualizarUsuario(ctx context.Context, input model.UsuarioInput) (model.Resultado, error) {
	resultado, err := controller.ActualizarUsuario(input)
	if err != nil {
		return resultado, err
	}

	return resultado, nil
}

// EliminarUsuario is the resolver for the eliminarUsuario field.
func (r *mutationResolver) EliminarUsuario(ctx context.Context, input model.UsuarioInput) (model.Resultado, error) {
	resultado, err := controller.EliminarUsuario(input)
	if err != nil {
		return resultado, err
	}

	return resultado, nil
}

// AutentificarUsuario is the resolver for the autentificarUsuario field.
func (r *queryResolver) AutentificarUsuario(ctx context.Context, input model.UsuarioInput) (model.TokenAuth, error) {
	tokenAuth, err := controller.AutentificarUsuario(input)
	if err != nil {
		tokenAuth = model.TokenAuth{Exito: false}
		return tokenAuth, err
	}
	return tokenAuth, nil
}

// TomarUsuarioXRecursos is the resolver for the tomarUsuarioXRecursos field.
func (r *queryResolver) TomarUsuarioXRecursos(ctx context.Context, id int64) (model.UsuarioXRecursos, error) {
	resultado, err := controller.TomarUsuarioXRecursos(id)
	if err != nil {
		resultado = model.UsuarioXRecursos{}
		return resultado, err
	}
	return resultado, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
/*
	type Resolver struct{}
*/
