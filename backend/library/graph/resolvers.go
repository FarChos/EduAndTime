package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"libreria/controller"
	"libreria/graph/model"
)

type Resolver struct{}

// ActualizarNumDescargas is the resolver for the actualizarNumDescargas field.
func (r *mutationResolver) ActualizarNumDescargas(ctx context.Context, id int) (model.Resultado, error) {
	resultado, err := controller.ActualizarNumDescargas(id)
	if err != nil {
		return resultado, err
	}
	return resultado, nil
}

// CalificarRecurso is the resolver for the calificarRecurso field.
func (r *mutationResolver) CalificarRecurso(ctx context.Context, id int, calificacion float64) (model.Resultado, error) {
	resultado, err := controller.CalificarRecurso(id, calificacion)
	if err != nil {
		return resultado, err
	}
	return resultado, nil
}

// SubirRecurso is the resolver for the subirRecurso field.
func (r *mutationResolver) SubirRecurso(ctx context.Context, input model.RecursoInput) (model.Resultado, error) {
	resultado, err := controller.SubirRecurso(input)
	if err != nil {
		return resultado, err
	}
	return resultado, nil
}

// EliminarRecurso is the resolver for the eliminarRecurso field.
func (r *mutationResolver) EliminarRecurso(ctx context.Context, id int) (model.Resultado, error) {
	resultado, err := controller.EliminarRecurso(id)
	if err != nil {
		return resultado, err
	}
	return resultado, nil
}

// BuscarRecursos is the resolver for the BuscarRecursos field.
func (r *queryResolver) BuscarRecursos(ctx context.Context, input model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	recursos, err := controller.BuscarRecursos(input)
	if err != nil {
		return nil, err
	}
	return recursos, nil
}

// TomarRecursos is the resolver for the tomarRecursos field.
func (r *queryResolver) TomarRecursos(ctx context.Context, ides []*int) ([]*model.RecursoMuestra, error) {
	recursos, err := controller.TomarRecursos(ides)
	if err != nil {
		return nil, err
	}
	return recursos, nil
}

// TomarRecurso is the resolver for the tomarRecurso field.
func (r *queryResolver) TomarRecurso(ctx context.Context, id int) (model.Recurso, error) {
	recurso, err := controller.TomarRecurso(id)
	if err != nil {
		return recurso, err
	}
	return recurso, nil
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
func (r *queryResolver) BuscarRecursosIniciales(ctx context.Context, input model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	recursos, err := controller.BuscarRecursosIniciales(input)
	if err != nil {
		return nil, err
	}
	return recursos, nil
}
*/
