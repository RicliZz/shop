package products

import (
	"fmt"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutesProduct(router *mux.Router) {
	router.HandleFunc("/products", h.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", h.GetProductById).Methods("GET")
}

func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("missing product with this ID"))
		return
	}
	productID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("error getting product by id: %w", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
	}
	utils.WriteJSON(w, http.StatusOK, products)
}
