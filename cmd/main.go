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
		return
	}
	initStorage(database)
	server := api.NewAPIServer(":8080", database)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to database")
}
