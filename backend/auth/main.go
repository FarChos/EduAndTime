package main

import (
	"ModuloAuth/db"
	"ModuloAuth/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
)

const defaultPort = "8080"

func main() {
	db.InitMariaDB("root:tearsthemoon@tcp(mariadb:3306)/EduAndTime?charset=utf8mb4&parseTime=True&loc=Local")
	db.InitCouchbase("couchbase://couchbase", "administrador", "tearsthemoon")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Configuración de CORS
	corsOptions := handlers.AllowedOrigins([]string{"http://localhost:5173"}) // Permitir solo este origen
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// Envolver srv con CORS
	http.Handle("/query", handlers.CORS(corsOptions, corsMethods, corsHeaders)(srv))

	log.Printf("Conéctate a http://localhost:%s/ para acceder al GraphQL playground", port)
	defer db.CloseMariaDB()
	defer db.CloseCouchbase()

	// Iniciar el servidor
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
