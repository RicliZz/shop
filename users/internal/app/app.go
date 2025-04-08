package app

import (
	"github.com/RicliZz/shop/users/internal/handlers"
	serviceStorage "github.com/RicliZz/shop/users/internal/services/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(addr string, router *gin.Engine) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           addr,
			Handler:        router.Handler(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	router := gin.Default()
	api := router.Group("/api/v1/")

	storageServ := serviceStorage.NewGetFullItemsInStorage()
	storageHandler := handlers.NewStorageHandler(storageServ)

	storageHandler.InitRoutes(api)

	serv := NewServer(os.Getenv("AddrServ"), router)
	serv.httpServer.ListenAndServe()
}
