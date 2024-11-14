package models

import (
	"log"

	"github.com/couchbase/gocb/v2"
)

var CouchbaseCluster *gocb.Cluster

// InitCouchbase inicializa la conexión a la base de datos Couchbase.
func InitCouchbase(connStr, username, password string) {
	var err error
	CouchbaseCluster, err = gocb.Connect(connStr, gocb.ClusterOptions{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Fatalf("Error al conectar con Couchbase: %v", err)
	}

	if _, err := CouchbaseCluster.Ping(nil); err != nil {
		log.Fatalf("Error al hacer ping a Couchbase: %v", err)
	}
	log.Println("Conexión a Couchbase establecida con éxito")
}

func CloseCouchbase() {
	if CouchbaseCluster != nil {
		if err := CouchbaseCluster.Close(nil); err != nil {
			log.Printf("Error al cerrar la conexión con Couchbase: %v", err)
		}
	}
}
