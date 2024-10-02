package api

import (
	"database/sql"
	_ "github.com/RiCliZz/shop/docs"
	"github.com/RiCliZz/shop/services/address"
	"github.com/RiCliZz/shop/services/cart"
	"github.com/RiCliZz/shop/services/orders"
	"github.com/RiCliZz/shop/services/products"
	"github.com/RiCliZz/shop/services/users"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/http-swagger/v2"
	"log"
	"net/http"
)

type APIServer struct {
	addr  string
	db    *sql.DB
	redis *redis.Client
}

func NewAPIServer(addr string, db *sql.DB, rdb *redis.Client) *APIServer {
	return &APIServer{addr: addr, db: db, redis: rdb}
}

func (s *APIServer) Start() error {
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The URL pointing to API definition
	))
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	addressStore := address.NewStore(s.db)
	userStore := users.NewStore(s.db)
	productStore := products.NewStore(s.db)
	cartStore := cart.NewStore(s.db)
	orderStore := orders.NewStore(s.db)

	userHandler := users.NewHandler(userStore, addressStore)
	userHandler.RegisterRoutesUser(subrouter)

	productHandler := products.NewHandler(productStore, userStore, s.redis)
	productHandler.RegisterRoutesProduct(subrouter)

	addressHandler := address.NewHandler(userStore, addressStore)
	addressHandler.RegisterRouterAddresses(subrouter)

	cartHandler := cart.NewHandler(cartStore, userStore, productStore)
	cartHandler.RegisterRoutesCart(subrouter)

	ordersHandler := orders.NewHandler(userStore, cartStore, orderStore, productStore)
	ordersHandler.RegisterRoutesOrders(subrouter)

	log.Println("Listening on " + s.addr)
	return http.ListenAndServe(s.addr, router)
}
