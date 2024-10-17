package main

import (
	"fmt"
	"log"

	"github.com/couchbase/gocb/v2"
)

func main() {
	// Conéctate al clúster
	cluster, err := gocb.Connect("couchbase://127.0.0.1", gocb.ClusterOptions{
		Username: "administrador",
		Password: "tearsthemoon",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Abre el bucket
	bucket := cluster.Bucket("EduAndTime")

	// Obtén una colección para trabajar
	collection := bucket.DefaultCollection()

	// Insertar un documento
	doc := map[string]string{
		"type":   "usuario",
		"nombre": "Juan Perez",
		"correo": "juan.perez@example.com",
	}
	_, err = collection.Insert("user_12345", doc, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documento insertado!")
}
