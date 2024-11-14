package main

import (
	"libreria/graph"
	"libreria/models"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
)

const defaultPort = "8081"

func main() {
	models.InitMariaDB("root:tearsthemoon@tcp(mariadb:3306)/EduAndTime?charset=utf8mb4&parseTime=True&loc=Local")
	models.InitCouchbase("couchbase://couchbase", "administrador", "tearsthemoon")

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	corsOptions := handlers.AllowedOrigins([]string{"http://localhost:5173"}) // Permitir solo este origen
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})
	corsHeaders := handlers.AllowedHeaders([]string{"Content-Type"})

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", handlers.CORS(corsOptions, corsMethods, corsHeaders)(srv))
	//http.Handle("/query", srv)

	log.Printf("Conéctate a http://localhost:%s/ para acceder al GraphQL playground", port)
	defer models.CloseMariaDB()
	defer models.CloseCouchbase()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
