package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var MariaDB *sql.DB

func InitMariaDB(dsn string) {
	var err error
	MariaDB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al conectar con MariaDB: %v", err)
	}

	if err = MariaDB.Ping(); err != nil {
		log.Fatalf("Error al hacer ping a la base de datos: %v", err)
	}
}

func CloseMariaDB() {
	if MariaDB != nil {
		if err := MariaDB.Close(); err != nil {
			log.Printf("Error al cerrar la conexión con MariaDB: %v", err)
		}
	}
}
