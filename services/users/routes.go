package users

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/regiter", registerHandler).Methods("POST")
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

}
