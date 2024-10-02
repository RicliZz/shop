package cart

import (
	"github.com/RiCliZz/shop/responses"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store        types.CartStore
	userStore    types.UserStore
	productStore types.ProductStore
}

func NewHandler(store types.CartStore, userStore types.UserStore, productStore types.ProductStore) *Handler {
	return &Handler{
		store:        store,
		userStore:    userStore,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutesCart(router *mux.Router) {
	router.HandleFunc("/cart", auth.WithJWTAuth(h.AddToCartHandler, h.userStore)).Methods("POST")
	router.HandleFunc("/cart", auth.WithJWTAuth(h.checkCart, h.userStore)).Methods("GET")
	router.HandleFunc("/cart", auth.WithJWTAuth(h.clearCartHandler, h.userStore)).Methods("DELETE")

}

// @Summary full clear you're cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} responses.SuccessResponse "Del all products in you're cart"
// @Router /cart [DELETE]
func (h *Handler) clearCartHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Not authorized",
		})
	}
	if err := h.store.DeleteCart(userId); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Failed to clear cart",
			Details: err.Error(),
		})
	}
	utils.WriteJSON(w, http.StatusOK, responses.SuccessResponse{
		Success: true,
		Data:    "You have successfully cleared cart",
	})
}

// @Summary get cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} types.CartItem "All products in cart"
// @Router /cart [GET]
func (h *Handler) checkCart(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Not auth"})
		return
	}
	products, err := h.store.CheckCart(userId)
	if err != nil {
		utils.ErrorJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Not auth",
		})
	}
	utils.WriteJSON(w, http.StatusOK, products)
}

// @Summary add product in cart
// @Tags Cart
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param product body types.AddToCart true "Product Id, quantity"
// @Success 200 {object} responses.SuccessResponse "success add"
// @Router /cart [POST]
func (h *Handler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	var product types.AddToCart
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}
	_, err := h.productStore.GetProductByID(product.Id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "product not found",
			Details: err.Error(),
		})
		return
	}
	err = h.store.AddToCart(userId, product.Id, product.Quantity)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "failed to add to cart",
			Details: err.Error(),
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, responses.SuccessResponse{
		Success: true,
	})
}
