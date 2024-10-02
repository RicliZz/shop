package orders

import (
	"fmt"
	"github.com/RiCliZz/shop/responses"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	userStore    types.UserStore
	cartStore    types.CartStore
	store        types.OrderStore
	productStore types.ProductStore
}

func NewHandler(userStore types.UserStore, cartStore types.CartStore, store types.OrderStore, productStore types.ProductStore) *Handler {
	return &Handler{
		userStore:    userStore,
		cartStore:    cartStore,
		store:        store,
		productStore: productStore,
	}
}

func (h *Handler) RegisterRoutesOrders(router *mux.Router) {
	router.HandleFunc("/order", auth.WithJWTAuth(h.newOrderHandler, h.userStore)).Methods("POST")
}

// @Summary Create order
// @Tags Order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {int} int "Success order"
// @Router /order [POST]
func (h *Handler) newOrderHandler(w http.ResponseWriter, r *http.Request) {
	var total float64
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusUnauthorized, responses.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}
	products, err := h.cartStore.CheckCart(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal server error",
			Details: err.Error(),
		})
		return
	}
	for _, v := range products {
		pr, err := h.productStore.GetProductByName(v.Name)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
				Error:   "Internal server error",
				Details: err.Error(),
			})
			return
		}
		pr.Quantity -= v.Quantity
		if pr.Quantity <= 0 {
			utils.WriteJSON(w, http.StatusBadRequest, responses.ErrorResponse{
				Error:   "Item end",
				Details: v.Name + fmt.Sprintf(" %d in storage", pr.Quantity+v.Quantity),
			})
			return
		}
		err = h.productStore.UpdateProduct(pr)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
				Error:   "Internal server error",
				Details: err.Error(),
			})
			return
		}
		total += v.Price
	}
	idOrder, err := h.store.CreateNewOrder(id, total)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "Internal server error",
			Details: err.Error(),
		})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, idOrder)
}
