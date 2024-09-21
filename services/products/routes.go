package products

import (
	"fmt"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
}

func NewHandler(store types.ProductStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutesProduct(router *mux.Router) {
	router.HandleFunc("/products", h.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", h.GetProductById).Methods("GET")

	//ONLY FOR ADMIN
	router.HandleFunc("/product/create", auth.WithJWTAuth(h.CreateProduct, h.userStore)).Methods("POST")
}

// Продукт по ID
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

// Вся продукция
func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, product); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	if err := utils.Validator.Struct(product); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}
	if err := h.store.CreateProduct(&product); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("error creating product: %w", err))
	}
	utils.WriteJSON(w, http.StatusCreated, product)
}
