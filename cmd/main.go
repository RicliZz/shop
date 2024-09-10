package main

import (
	"database/sql"
	"github.com/RiCliZz/shop/cmd/api"
	"github.com/RiCliZz/shop/db"
	"log"
)

func main() {
	database, err := db.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	initDB(database)
	server := api.NewAPIServer(":8080", database)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}

}

func initDB(db *sql.DB) {
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

}
