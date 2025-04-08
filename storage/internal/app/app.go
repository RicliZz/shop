package app

import (
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
			Handler:        router,
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
	//api := router.Group("/api/v1/storage")

	serv := NewServer(os.Getenv("AddrServ"), router)
	serv.httpServer.ListenAndServe()
}
