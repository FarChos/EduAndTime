package main

import (
	"libreria/graph"
	"libreria/middleware"
	"libreria/models"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const defaultPort = "8081"

func main() {
	// Inicializar bases de datos
	models.InitMariaDB("root:tearsthemoon@tcp(mariadb:3306)/EduAndTime?charset=utf8mb4&parseTime=True&loc=Local")
	models.InitCouchbase("couchbase://couchbase", "administrador", "tearsthemoon")
	defer models.CloseMariaDB()
	defer models.CloseCouchbase()

	// Configuración del puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Crear instancia de Echo
	e := echo.New()

	// Middlewares globales
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderContentType, "Authorization"},
	}))

	// Configurar el servidor GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// Ruta para el playground
	e.GET("/", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))

	// Rutas protegidas con validación de token
	e.POST("/query", middleware.ValidateTokenMiddleware(echo.WrapHandler(srv)))

	// Iniciar el servidor
	log.Printf("Conéctate a http://localhost:%s/ para acceder al GraphQL playground", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
