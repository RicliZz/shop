package api

import (
	"database/sql"
	_ "github.com/RiCliZz/shop/docs"
	"github.com/RiCliZz/shop/services/address"
	"github.com/RiCliZz/shop/services/products"
	"github.com/RiCliZz/shop/services/users"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{addr: addr, db: db}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The URL pointing to API definition
	))
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	addressStore := address.NewStore(s.db)
	userStore := users.NewStore(s.db)
	userHandler := users.NewHandler(userStore, addressStore)
	userHandler.RegisterRoutesUser(subrouter)

	productStore := products.NewStore(s.db)
	productHandler := products.NewHandler(productStore, userStore)
	productHandler.RegisterRoutesProduct(subrouter)

	addressHandler := address.NewHandler(userStore, addressStore)
	addressHandler.RegisterRouterAddresses(subrouter)

	log.Println("Listening on " + s.addr)
	return http.ListenAndServe(s.addr, router)
}
