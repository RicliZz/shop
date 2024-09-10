package db

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewPostgres() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database, err := sql.Open("postgres", os.Getenv("PSQL"))
	if err != nil {
		return nil, err
	}
	return database, nil
}
