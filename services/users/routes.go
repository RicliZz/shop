package users

import (
	"fmt"
	"github.com/RiCliZz/shop/services/auth"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "gopkg.in/gomail.v2"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutesUser(router *mux.Router) {
	router.HandleFunc("/register", h.registerHandler).Methods("POST")
	router.HandleFunc("/login", h.loginHandler).Methods("POST")
	router.HandleFunc("/confirm", h.confirmEmailHandler).Methods("GET")
	router.HandleFunc("/profile", auth.WithJWTAuth(h.profileHandler, h.store)).Methods("GET")

	//only for admin
	router.HandleFunc("/user/{id}", auth.WithJWTAdminAuth(h.handleGetUser, h.store)).Methods("GET")
}
func (h *Handler) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	err := h.store.CheckToken(r.URL.Query().Get("token"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
	}
	utils.WriteJSON(w, http.StatusOK, "SUCCESS")
}

func (h *Handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("not authorized!!!"))
		return
	}
	u, err := h.store.GetUserByIDForProfile(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("user not found"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, u)
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
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	if !auth.ComparePass(u.Password, user.Password) {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("invalid password"))
		return
	}
	if !u.Email_verified {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("email not verified"))
		return
	}
	godotenv.Load()
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := auth.CreateJWT(secret, u.Id, u.Role)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
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
	token, err := utils.GenerateUUID()
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	//Если всё файн - регистрируем
	err = h.store.CreateAcc(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
		Password:  hashedPass,
	})
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, "Please check ur email")

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, fmt.Errorf("id not found"))
		return
	}
	id, err := strconv.Atoi(str)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err)
		return
	}
	user, err := h.store.GetUserByID(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
