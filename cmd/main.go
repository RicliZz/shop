package main

import (
	"context"
	"database/sql"
	"github.com/RiCliZz/shop/cmd/api"
	"github.com/RiCliZz/shop/db"
	"github.com/redis/go-redis/v9"
	"log"
)

// @title InternetShop
// @version 1.0
// @description API Server for Internet Shop

// @host localhost:8080
// @BasePath /api/v1/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	database, err := db.NewPostgres()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatal(err)
	}
	err = initDB(database)
	if err != nil {
		log.Fatal(err)
	}
	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to redis server")
	server := api.NewAPIServer(":8080", database, rdb)
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func initDB(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	log.Println("Successfully connected to database")
	return nil

}
