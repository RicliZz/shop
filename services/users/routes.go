package users

import (
	"fmt"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutesUser(router *mux.Router) {
	router.HandleFunc("/register", h.registerHandler).Methods("POST")
	router.HandleFunc("/login", h.loginHandler).Methods("POST")
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	//Приведение к JSON
	var user types.UserLoginPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	//Проверка валидатором
	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("invalid user: %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("user with email %v not found", user.Email))
		return
	}
	if !auth.ComparePass(u.Password, user.Password) {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, "Welcome!")
}

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	//Приведение к JSON
	var user types.UserRegisterPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	//Проверка валидатором
	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("invalid user: %v", errors))
		return
	}
	//Проверка на уже существующий email
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("user with this email already exists"))
		return
	}

	//Хэшируем пароль
	hashedPass, err := auth.HashPass(user.Password)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
	}
	//Если всё файн - регистрируем
	err = h.store.CreateAcc(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPass,
	})
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}
