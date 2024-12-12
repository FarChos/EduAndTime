package services

import (
	"fmt"
	"libreria/graph/model"
	"strconv"
	"strings"
)

func placeholders(n int) string {
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

func TomarIdes(recursos []*model.RecursoMuestra) []int {
	// Manejar el caso en que el slice sea nulo o vacío
	if recursos == nil {
		return []int{}
	}

	// Crear un slice para los IDs
	ides := make([]int, len(recursos))
	for i, recurso := range recursos {
		ides[i] = recurso.ID
	}
	return ides
}

func ManejoEtiquetas(etiquetasRaw []interface{}) ([]*string, error) {
	if len(etiquetasRaw) == 0 {
		// Si no hay etiquetas, se devuelve nil, sin error
		return nil, nil
	}

	// Si hay etiquetas, se procesan
	etiquetas := make([]*string, len(etiquetasRaw))
	for i, etiqueta := range etiquetasRaw {
		etiquetaStr, ok := etiqueta.(string)
		if !ok {
			return nil, fmt.Errorf("una de las etiquetas no es un string")
		}
		etiquetas[i] = &etiquetaStr
	}
	return etiquetas, nil
}

// ! ProcesarResultadoCouchbase
func ProcesarResultadoCouchbase(resultados map[string]map[string]interface{}) ([]*model.RecursoMuestra, error) {
	var recursosMuestra []*model.RecursoMuestra
	if resultados == nil {
		return recursosMuestra, fmt.Errorf("Couchbase no  devolvió resultados a procesar: %w")
	}

	for _, resultado := range resultados {
		var recursoMuestra model.RecursoMuestra

		// Procesar etiquetas
		if etiquetasRaw, ok := resultado["etiquetas"].([]interface{}); ok {
			etiquetas, err := ManejoEtiquetas(etiquetasRaw)
			if err != nil {
				return recursosMuestra, fmt.Errorf("error al procesar las etiquetas: %w", err)
			}
			recursoMuestra.Etiquetas = etiquetas
		}

		// Procesar ID del documento
		if docID, ok := resultado["id"].(string); ok {
			// Dividir el string usando "::" como delimitador
			parts := strings.Split(docID, "::")

			if len(parts) > 1 {
				// Tomar la segunda parte (el número) y convertirla a int
				intID, err := strconv.Atoi(parts[1])
				if err != nil {
					return recursosMuestra, fmt.Errorf("error al convertir el ID del documento a int: %w", err)
				}
				recursoMuestra.ID = intID
			} else {
				return recursosMuestra, fmt.Errorf("error al procesar el ID del documento: formato incorrecto")
			}
		} else {
			return recursosMuestra, fmt.Errorf("error al procesar el ID del documento: falta o es de tipo incorrecto")
		}

		// Procesar calificación
		if calificacionRaw, ok := resultado["calificacion"].(float64); ok {
			recursoMuestra.Calificacion = &calificacionRaw
		}

		// Agregar recurso procesado a la lista
		recursosMuestra = append(recursosMuestra, &recursoMuestra)
	}

	return recursosMuestra, nil
}

func UnificarRecursosMuestra(
	recursosMuestramariadb []*model.RecursoMuestra,
	recursosMuestraCouchbase []*model.RecursoMuestra,
) ([]*model.RecursoMuestra, error) {
	var recursosMuestraFinales []*model.RecursoMuestra

	// Verificar si ambos resultados son nil o vacíos
	if len(recursosMuestramariadb) == 0 && len(recursosMuestraCouchbase) == 0 {
		return nil, fmt.Errorf("Couchbase y MariaDB no devolvieron resultados a procesar")
	}

	// Crear un mapa para indexar los resultados de Couchbase por ID
	couchbaseMap := make(map[int]*model.RecursoMuestra)
	for _, recurso := range recursosMuestraCouchbase {
		couchbaseMap[recurso.ID] = recurso
	}

	// Unir resultados, pero incluir solo los que están en Couchbase
	for _, recursoMariaDB := range recursosMuestramariadb {
		if recursoCouchbase, ok := couchbaseMap[recursoMariaDB.ID]; ok {
			// Combinar datos de Couchbase con MariaDB
			recursoMariaDB.Etiquetas = recursoCouchbase.Etiquetas
			recursoMariaDB.Calificacion = recursoCouchbase.Calificacion
			recursosMuestraFinales = append(recursosMuestraFinales, recursoMariaDB)
		}
	}

	return recursosMuestraFinales, nil
}
func UnificarRecursosMuestraMariadb(recursosAleatorios []*model.RecursoMuestra, recursosUltimos []*model.RecursoMuestra) ([]*model.RecursoMuestra, error) {
	// Validar si ambos resultados están vacíos
	if len(recursosAleatorios) == 0 && len(recursosUltimos) == 0 {
		return nil, fmt.Errorf("Mariadb no devolvió recursos para unificar")
	}

	// Crear un nuevo slice para almacenar los recursos combinados
	var recursosUnificados []*model.RecursoMuestra

	// Agregar recursos aleatorios si están presentes
	if len(recursosAleatorios) > 0 {
		recursosUnificados = append(recursosUnificados, recursosAleatorios...)
	}

	// Agregar recursos últimos si están presentes
	if len(recursosUltimos) > 0 {
		recursosUnificados = append(recursosUnificados, recursosUltimos...)
	}

	return recursosUnificados, nil
}

func TomarIdesDePunteros(ides []*int) ([]int, error) {
	var docIdes []int
	if len(ides) == 0 {
		return docIdes, fmt.Errorf("la lista de ids esta vacia")
	}
	for _, id := range ides {
		docIdes = append(docIdes, *id)
	}
	return docIdes, nil
}
