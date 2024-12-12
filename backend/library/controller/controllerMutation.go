package controller

import (
	"fmt"
	"libreria/graph/model"
	"libreria/services"
)

// ! actualiza el numero de descargas en couchbase
func ActualizarNumDescargas(id int) (model.Resultado, error) {
	resultado := model.Resultado{
		Exito: false,
	}
	resultadoRows, err := services.TomarNoDescargasEnCouchBase(id)
	if err != nil {
		return resultado, err
	}
	mapaResultado, err := services.DevolverResultadoClaveValor(resultadoRows)
	if err != nil {
		return resultado, err
	}
	noNewDescargas, err := services.DevolverNuevoNoDescargas(mapaResultado)
	if err != nil {
		return resultado, err
	}
	err = services.ActualizarNoDescargas(noNewDescargas, id)
	if err != nil {
		return resultado, err
	}
	resultado.Exito = true
	return resultado, nil
}

// ! cambia la calificacion promedio del recurso en couchbase
func CalificarRecurso(id int, calificacion float64) (model.Resultado, error) {
	resultado := model.Resultado{
		Exito: false,
	}
	resultQuery, err := services.ObtenerCalificaciones(id)
	if err != nil {
		return resultado, err
	}
	claveValor, err := services.DevolverResultadoClaveValor(resultQuery)
	if err != nil {
		return resultado, err
	}
	nuevasCalificaciones, err := services.DevolverNuevaCalificacion(claveValor, calificacion)
	if err != nil {
		return resultado, err
	}
	err = services.ActualizarCalificacion(nuevasCalificaciones, id)
	if err != nil {
		return resultado, err
	}
	return resultado, nil
}
func SubirRecurso(input model.RecursoInput) (model.Resultado, error) {
	// Inicializar el resultado
	resultado := model.Resultado{
		Exito: false,
	}
	nombreRecurso, err := services.GuardarRecurso(input, "/recursos")
	if err != nil {
		return resultado, fmt.Errorf("error al crear el recurso en MariaDB: %w", err)
	}
	// Crear el recurso en MariaDB
	idLibroNuevo, err := services.CrearRecursoEnMariadb(input, nombreRecurso)
	if err != nil {
		return resultado, fmt.Errorf("error al crear el recurso en MariaDB: %w", err)
	}

	// Convertir etiquetas
	etiquetas, err := services.ConvertirEtiquetas(input.Etiquetas)
	if err != nil {
		// Considerar eliminar el recurso recién creado en MariaDB
		_ = services.EliminarRecursoEnMariadb(idLibroNuevo) // Manejo silencioso del error
		return resultado, fmt.Errorf("error al convertir etiquetas: %w", err)
	}

	// Crear el recurso en Couchbase
	err = services.CrearRecursoEnCouchbase(idLibroNuevo, etiquetas)
	if err != nil {
		// Revertir cambios en MariaDB si Couchbase falla
		_ = services.EliminarRecursoEnMariadb(idLibroNuevo) // Manejo silencioso del error
		return resultado, fmt.Errorf("error al crear el recurso en Couchbase: %w", err)
	}

	// Actualizar el resultado en caso de éxito
	resultado.Exito = true
	return resultado, nil
}

func EliminarRecurso(id int) (model.Resultado, error) {
	resultado := model.Resultado{
		Exito: false,
	}

	// Eliminar recurso en MariaDB
	if err := services.EliminarRecursoEnMariadb(int64(id)); err != nil {
		return resultado, fmt.Errorf("error al eliminar recurso en MariaDB: %w", err)
	}

	// Eliminar recurso en Couchbase
	if err := services.EliminarDeCouchbase(int64(id)); err != nil {
		// Considerar realizar un rollback en MariaDB si es necesario
		return resultado, fmt.Errorf("error al eliminar recurso en Couchbase: %w", err)
	}

	// Marcar como éxito
	resultado.Exito = true
	return resultado, nil
}
