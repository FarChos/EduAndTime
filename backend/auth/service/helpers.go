package service

import (
	"ModuloAuth/graph/model"
	"fmt"
	"strconv"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

const maxPasswordLength = 72 // Tamaño máximo recomendado por bcrypt

// ! HashContrasena genera un hash seguro a partir de una contraseña usando bcrypt.
func HashContrasena(contrasena string) ([]byte, error) {
	// Verificar longitud de la contraseña
	if len(contrasena) > maxPasswordLength {
		return nil, fmt.Errorf("la contraseña excede el tamaño máximo permitido de %d caracteres", maxPasswordLength)
	}

	// Generar el hash de la contraseña
	contrasenaHash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error al hashear la contraseña: %v", err)
	}

	return contrasenaHash, nil
}

func GenerarToken(nombre string) (string, error) {
	horaExpiracion := time.Now().Add(1 * time.Hour).Unix()
	claims := &Claims{
		UserID: nombre,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: horaExpiracion,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("viajeAntesQueDestino"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
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

// ! DevolverUsuarioXRecursos procesa el mapa de strings y lo devuelve en forma de un tipo UsuarioXRecursos
func DevolverUsuarioXRecursos(ClaveValor map[string]interface{}) (model.UsuarioXRecursos, error) {
	var DatosRecursoXUsuario model.UsuarioXRecursos

	// Procesar la lista de favoritos
	if favoritos, ok := ClaveValor["docFavoritos"].([]interface{}); ok {
		var listaFavoritos []*int64
		for _, fav := range favoritos {
			if valor, ok := fav.(float64); ok { // JSON decodifica números como float64
				valorInt := int64(valor)
				listaFavoritos = append(listaFavoritos, &valorInt)
			}
		}
		DatosRecursoXUsuario.IdesFavoritos = listaFavoritos
	}

	// Procesar la lista de originados
	if originados, ok := ClaveValor["docOriginados"].([]interface{}); ok {
		var listaOriginados []*int64
		for _, orig := range originados {
			if valor, ok := orig.(float64); ok { // JSON decodifica números como float64
				valorInt := int64(valor)
				listaOriginados = append(listaOriginados, &valorInt)
			}
		}
		DatosRecursoXUsuario.IdesMisRecursos = listaOriginados
	}

	// Procesar el mapa de calificados
	if calificados, ok := ClaveValor["docCalificados"].(map[string]interface{}); ok {
		var listaCalificados []*model.RecursoCalificado
		for clave, valor := range calificados {
			// Convertir la clave a int64 y el valor a int32
			claveInt, err := strconv.ParseInt(clave, 10, 64)
			if err != nil {
				return DatosRecursoXUsuario, fmt.Errorf("error al convertir clave '%s' a int64: %w", clave, err)
			}
			if valorFloat, ok := valor.(float64); ok {
				valorInt := float64(valorFloat)
				listaCalificados = append(listaCalificados, &model.RecursoCalificado{
					ID:           claveInt,
					Calificacion: valorInt,
				})
			}
		}
		DatosRecursoXUsuario.RecursosCalificados = listaCalificados
	}

	return DatosRecursoXUsuario, nil
}
