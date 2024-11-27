package graph

// THIS CODE WILL BE UPDATED WITH SCHEMA CHANGES. PREVIOUS IMPLEMENTATION FOR SCHEMA CHANGES WILL BE KEPT IN THE COMMENT SECTION. IMPLEMENTATION FOR UNCHANGED SCHEMA WILL BE KEPT.

import (
	"context"
	"libreria/graph/model"
	"libreria/services"
)

type Resolver struct{}

// ActualizarNumDescargas is the resolver for the actualizarNumDescargas field.
func (r *mutationResolver) ActualizarNumDescargas(ctx context.Context, id int) (bool, error) {
	seActualizo, err := services.ActualizarNumDescargas(id)
	if err != nil {
		return false, err
	}
	return seActualizo, nil
}

// CalificarRecurso is the resolver for the calificarRecurso field.
func (r *mutationResolver) CalificarRecurso(ctx context.Context, id int, calificacion float64) (bool, error) {
	seCalifico, err := services.CalificarRecurso(id, float32(calificacion))
	if err != nil {
		return false, err
	}
	return seCalifico, nil
}

// SubirRecurso is the resolver for the subirRecurso field.
func (r *mutationResolver) SubirRecurso(ctx context.Context, input model.RecursoInput) (bool, error) {
	seSubio, err := services.CrearRecurso(input)
	if err != nil {
		return false, err
	}
	return seSubio, nil
}

// EliminarRecurso is the resolver for the eliminarRecurso field.
func (r *mutationResolver) EliminarRecurso(ctx context.Context, id int) (bool, error) {
	seElimino, err := services.EliminarRecurso(id)
	if err != nil {
		return false, err
	}
	return seElimino, nil
}

// BuscarRecursos is the resolver for the BuscarRecursos field.
func (r *queryResolver) BuscarRecursos(ctx context.Context, input model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	recursos, err := services.BuscarRecursos(input)
	if err != nil {
		return nil, err
	}
	return recursos, nil
}

// TomarRecurso is the resolver for the tomarRecurso field.
func (r *queryResolver) TomarRecurso(ctx context.Context, id int) (model.Recurso, error) {
	exito, err := services.TomarRecurso(id)
	if err != nil {
		return exito, err
	}
	return exito, nil
}

// UltimosRecursos is the resolver for the ultimosRecursos field.
func (r *queryResolver) UltimosRecursos(ctx context.Context) ([]*model.RecursoMuestra, error) {
	recursosUltimos, err := services.UltimosRecursos()
	if err != nil {
		return nil, err
	}
	return recursosUltimos, nil
}

// RecursosAleatorios is the resolver for the recursosAleatorios field.
func (r *queryResolver) RecursosAleatorios(ctx context.Context) ([]*model.RecursoMuestra, error) {
	recursosAleatorios, err := services.RecursosAleatorios()
	if err != nil {
		return nil, err
	}
	return recursosAleatorios, nil
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
func (r *queryResolver) RecomendarRecursos(ctx context.Context) ([]*model.RecursoMuestra, error) {
	panic("not implemented")
}
*/
