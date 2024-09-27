package products

import (
	"github.com/RiCliZz/shop/responses"
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
	router.HandleFunc("/product/create", auth.WithJWTAdminAuth(h.CreateProduct, h.userStore)).Methods("POST")
}

// Продукт по ID
// @Summary one product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} types.CreateProductPayload "Get one product with ID"
// @Router /product/{id} [GET]
func (h *Handler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: "bad_request: product id missing",
		})
		return
	}
	productID, err := strconv.Atoi(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "invalid ID",
		})
		return
	}
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "not found product with this ID",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, product)
}

// Вся продукция
// @Summary all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} []types.ShortProducts "Get all products"
// @Router /products [GET]
func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.store.GetProducts()
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error with getting products",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, product); err != nil {
		utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "parsing request body error",
		})
		return
	}
	if err := utils.Validator.Struct(product); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error validating request body",
		})
		return
	}
	if err := h.store.CreateProduct(&product); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error creating product",
		})
	}
	utils.WriteJSON(w, http.StatusCreated, product)
}
