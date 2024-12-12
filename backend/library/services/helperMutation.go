package services

import (
	"fmt"

	"github.com/couchbase/gocb/v2"
)

// ! DevolverResultadoClaveValor procesa las filas devualtas por couchbase y devuelve un mapa string
func DevolverResultadoClaveValor(rows *gocb.QueryResult) (map[string]interface{}, error) {
	// Asegurarse de cerrar el objeto rows
	defer rows.Close()

	var resultado map[string]interface{}

	// Iterar sobre las filas
	if rows.Next() {
		// Procesar la fila actual
		err := rows.Row(&resultado)
		if err != nil {
			return nil, fmt.Errorf("error al leer la fila: %w", err)
		}

		// Verificar si hay más filas, lo que indicaría datos inesperados
		if rows.Next() {
			return nil, fmt.Errorf("se encontraron múltiples filas, se esperaba solo una")
		}

		// Retornar el resultado si salió bien
		return resultado, nil
	}

	// Verificar si ocurrió algún error durante la iteración
	if rows.Err() != nil {
		return nil, fmt.Errorf("error en la iteración de las filas: %w", rows.Err())
	}

	// Si no se encontraron filas
	return nil, fmt.Errorf("no se encontró el documento")
}

// ! devuelve el nuevo numero de descargas
func DevolverNuevoNoDescargas(resultado map[string]interface{}) (int, error) {
	noDescargasRaw, ok := resultado["noDescargas"]
	if !ok {
		return 0, fmt.Errorf("el campo 'noDescargas' no existe")
	}

	// Convertir a float64 y luego a int
	noDescargasFloat, ok := noDescargasRaw.(float64)
	if !ok {
		return 0, fmt.Errorf("el campo 'noDescargas' no es del tipo esperado en el documento")
	}
	noDescargas := int(noDescargasFloat) + 1
	return noDescargas, nil
}

// !
func DevolverNuevaCalificacion(resultado map[string]interface{}, calificacion float64) ([]float64, error) {
	// Manejar "noCalificaciones"
	noCalificacionesRaw, ok := resultado["noCalificaciones"].(float64)
	if !ok {
		// Si no existe o no es del tipo esperado, devolver un error
		return nil, fmt.Errorf("el valor de 'noCalificaciones' no es válido o no está presente")
	}

	// Manejar "calificacion"
	calificacionActualRaw, ok := resultado["calificacion"].(float64)
	if !ok {
		// Si no existe o no es del tipo esperado, devolver un error
		return nil, fmt.Errorf("el valor de 'calificacion' no es válido o no está presente")
	}

	// Convertir valores para cálculos
	noCalificaciones := int(noCalificacionesRaw) // Convertir a entero para conteo
	if noCalificaciones < 0 {
		return nil, fmt.Errorf("el valor de 'noCalificaciones' no puede ser negativo")
	}

	calificacionActual := calificacionActualRaw
	if calificacionActual < 0 || calificacionActual > 10 { // Rango ejemplo
		return nil, fmt.Errorf("la calificación actual está fuera del rango esperado")
	}

	// Calcular nueva calificación
	newCalificacion := (calificacionActual*float64(noCalificaciones) + calificacion) / float64(noCalificaciones+1)
	newNoCalificaciones := float64(noCalificaciones + 1)

	// Crear resultado final
	resultadoFinal := []float64{newCalificacion, newNoCalificaciones}
	return resultadoFinal, nil
}
func ConvertirEtiquetas(etiquetas []*string) ([]string, error) {
	// Verificar si la entrada es nil o vacía
	if etiquetas == nil || len(etiquetas) == 0 {
		return nil, nil
	}

	// Crear el slice de resultado con capacidad inicial
	resultado := make([]string, 0, len(etiquetas))

	// Convertir etiquetas no nulas
	for _, etiqueta := range etiquetas {
		if etiqueta != nil && *etiqueta != "" { // Validar que no sea nil ni vacía
			resultado = append(resultado, *etiqueta)
		}
	}

	return resultado, nil
}
