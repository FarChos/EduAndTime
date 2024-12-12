package controller

import (
	"fmt"
	"libreria/graph/model"
	"libreria/services"
)

func TomarRecurso(id int) (model.Recurso, error) {
	var recurso model.Recurso
	recurso, err := services.TomarDeMariadb(id)
	if err != nil {
		return recurso, err
	}
	recurso, err = services.TomarDeCouchbase(recurso)
	return recurso, nil
}

func BuscarRecursos(input model.ParametrosBusqueda) ([]*model.RecursoMuestra, error) {
	var (
		recursosMuestraFinales []*model.RecursoMuestra
		recursosAleatorios     []*model.RecursoMuestra
		recursosUltimos        []*model.RecursoMuestra
		recursosUnificados     []*model.RecursoMuestra
		recursosCouchbase      []*model.RecursoMuestra
		recursosBrutoCouchbase map[string]map[string]interface{}
	)

	// Validar y ajustar los parámetros de entrada
	if input.Categoria != nil && input.Categoria.String() == "general" {
		input.Categoria = nil
	}

	// Validar cantidad
	if input.Cantidad == nil || *input.Cantidad < 1 {
		return nil, fmt.Errorf("cantidad inválida: debe ser mayor a 0")
	}

	// Tomar muestra aleatoria de MariaDB
	var err error
	recursosAleatorios, err = services.TomarMuestraAleatoriaMariadb(input)
	if err != nil {
		return nil, fmt.Errorf("error al obtener muestra aleatoria de MariaDB: %w", err)
	}

	// Tomar muestra de los últimos recursos si aplica
	if *input.Cantidad > 3 {
		recursosUltimos, err = services.TomarMuestraUltimaMariadb(input)
		if err != nil {
			return nil, fmt.Errorf("error al obtener últimos recursos de MariaDB: %w", err)
		}
	}

	// Unificar recursos de MariaDB
	recursosUnificados, err = services.UnificarRecursosMuestraMariadb(recursosAleatorios, recursosUltimos)
	if err != nil {
		return nil, fmt.Errorf("error al unificar recursos de MariaDB: %w", err)
	}

	// Obtener los IDs de los recursos unificados
	ides := services.TomarIdes(recursosUnificados)
	if len(ides) == 0 {
		return nil, fmt.Errorf("no se encontraron IDs válidos en los recursos de MariaDB")
	}

	// Consultar Couchbase usando los IDs
	recursosBrutoCouchbase, err = services.TomarMuestraCouchbase(ides, &input)
	if err != nil {
		return nil, fmt.Errorf("error al obtener recursos de Couchbase: %w", err)
	}

	// Procesar los resultados de Couchbase
	recursosCouchbase, err = services.ProcesarResultadoCouchbase(recursosBrutoCouchbase)
	if err != nil {
		return nil, fmt.Errorf("error al procesar resultados de Couchbase: %w", err)
	}

	// Unificar los resultados finales
	recursosMuestraFinales, err = services.UnificarRecursosMuestra(recursosUnificados, recursosCouchbase)
	if err != nil {
		return nil, fmt.Errorf("error al unificar recursos finales: %w", err)
	}

	return recursosMuestraFinales, nil
}

func TomarRecursos(ides []*int) ([]*model.RecursoMuestra, error) {
	var recursosMuestraMariadb []*model.RecursoMuestra
	var recursosMuestraCouchbase []*model.RecursoMuestra
	var recursosMuestraFinales []*model.RecursoMuestra
	var recursosMuestraBrutosCouchbase map[string]map[string]interface{}

	var docIdes []int
	// Tomar los valores de los punteros
	docIdes, err := services.TomarIdesDePunteros(ides)
	if err != nil {
		return nil, fmt.Errorf("error al tomar los ides de los punteros: %w", err)
	}
	// Consultar  a mariadb usando los IDs
	recursosMuestraMariadb, err = services.TomarMuestraconIdes(docIdes)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los recursos por ides en MariaDB: %w", err)
	}
	// Consultar Couchbase usando los IDs
	recursosMuestraBrutosCouchbase, err = services.TomarMuestraCouchbase(docIdes, nil)
	if err != nil {
		return nil, fmt.Errorf("error al obtener recursos de Couchbase: %w", err)
	}
	// Procesar los resultados de Couchbase
	recursosMuestraCouchbase, err = services.ProcesarResultadoCouchbase(recursosMuestraBrutosCouchbase)
	if err != nil {
		return nil, fmt.Errorf("error al procesar resultados de Couchbase: %w", err)
	}
	// Unificar los resultados finales
	recursosMuestraFinales, err = services.UnificarRecursosMuestra(recursosMuestraMariadb, recursosMuestraCouchbase)
	if err != nil {
		return nil, fmt.Errorf("error al unificar recursos finales: %w", err)
	}

	return recursosMuestraFinales, nil
}
