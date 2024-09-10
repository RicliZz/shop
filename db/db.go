package db

import (
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewPostgres() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}
	db, err := sql.Open("postgres", os.Getenv("PSQL"))
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
