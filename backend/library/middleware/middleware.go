package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte("viajeAntesQueDestino") // Clave secreta compartida

// ValidateTokenMiddleware valida el token JWT
func ValidateTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Obtener el token desde los headers
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "No se proporcionó el token")
		}

		// Extraer el token del encabezado
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Si el token no sigue el formato esperado
			return echo.NewHTTPError(http.StatusUnauthorized, "Formato del token inválido")
		}

		// Parsear y validar el token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma no válido")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token inválido o expirado")
		}

		// Continuar con la siguiente función
		return next(c)
	}
}
