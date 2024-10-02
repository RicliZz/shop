package products

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RiCliZz/shop/responses"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	store     types.ProductStore
	userStore types.UserStore
	rdb       *redis.Client
}

func NewHandler(store types.ProductStore, userStore types.UserStore, rdb *redis.Client) *Handler {
	return &Handler{store: store, userStore: userStore, rdb: rdb}
}

func (h *Handler) RegisterRoutesProduct(router *mux.Router) {
	router.HandleFunc("/products", h.GetAllProducts).Methods("GET")
	router.HandleFunc("/product/{id}", h.GetProductById).Methods("GET")
	//ONLY FOR ADMIN
	router.HandleFunc("/product/create", auth.WithJWTAdminAuth(h.CreateProduct, h.userStore)).Methods("POST")
	router.HandleFunc("/product/add", auth.WithJWTAdminAuth(h.AddProductHandler, h.userStore)).Methods("POST")
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

	redisKey := fmt.Sprintf("productID:%d", productID)
	productData, err := h.rdb.Get(context.Background(), redisKey).Result()
	if err == nil && productData != "" {
		var product types.CreateProductPayload
		err = json.Unmarshal([]byte(productData), &product)
		if err != nil {
			utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "invalid product data",
			})
			return
		}
	}
	product, err := h.store.GetProductByID(productID)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "not found product with this ID",
		})
		return
	}

	productJSON, err := json.Marshal(product)
	if err == nil {
		h.rdb.Set(context.Background(), redisKey, string(productJSON), time.Minute*5)
	} else {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "internal server error",
		})
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
	cashedProducts, err := h.rdb.Get(context.Background(), "products").Result()
	if err == nil && cashedProducts != "" {
		var products []*types.ShortProducts
		err = json.Unmarshal([]byte(cashedProducts), &products)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "error with unmarshaling cached products",
			})
			return
		}
		// Отправляем кешированные данные
		utils.WriteJSON(w, http.StatusOK, products)
		log.Println("REDIS")
		return
	}

	products, err := h.store.GetProducts()
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error with getting products",
		})
		return
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error with marshaling products",
		})
		return
	}

	h.rdb.Set(context.Background(), "products", productsJSON, time.Minute*5)
	utils.WriteJSON(w, http.StatusOK, products)
}

// @Summary create new product
// @Description ONLY FOR ADMIN
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body types.CreateProductPayload true "Product data"
// @Success 200 {object} types.CreateProductPayload "Product"
// @Router /product/create [POST]
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product types.CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
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

// @Summary addQuantityProduct
// @Description ONLY FOR ADMIN
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body types.AddProduct true "Product data"
// @Success 200 {object} types.AddProduct "Product"
// @Router /product/add [POST]
func (h *Handler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var product types.AddProduct
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "parsing request body error",
		})
		return
	}

	prod, err := h.store.GetProductByName(product.Name)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error getting product",
		})
		return
	}
	log.Println(prod)
	prod.Quantity += product.Quantity
	if err = h.store.UpdateProduct(prod); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error adding product",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"name":     prod.Name,
		"quantity": prod.Quantity,
	})
}
