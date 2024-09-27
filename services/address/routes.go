package address

import (
	"github.com/RiCliZz/shop/responses"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	userStore    types.UserStore
	addressStore types.AddressStore
}

func NewHandler(userStore types.UserStore, addressStore types.AddressStore) *Handler {
	return &Handler{userStore, addressStore}
}

func (h *Handler) RegisterRouterAddresses(router *mux.Router) {
	router.HandleFunc("/ADDress", auth.WithJWTAuth(h.addAddressHandler, h.userStore)).Methods("POST")
}

// @Summary Create Address
// @Description Create address
// @Tags Address
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param address body types.AddressPayload true "Address data"
// @Success 201 {object} responses.SuccessResponse "User with this ID: "
// @Router /ADDress [Post]
func (h *Handler) addAddressHandler(w http.ResponseWriter, r *http.Request) {
	user_id, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
			Error: "Not Authorized",
		})
		return
	}
	var address types.AddressPayload
	if err := utils.ParseJSON(r, &address); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Error:   "Invalid Request",
			Details: err.Error(),
		})
		return
	}
	if err := h.addressStore.CreateNewAddress(user_id, types.Address{
		City:      address.City,
		Street:    address.Street,
		House:     address.House,
		Apartment: address.Apartment,
	}); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Failed to create new address",
		})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, responses.SuccessResponse{
		Success: true,
		Data:    address,
	})

}
