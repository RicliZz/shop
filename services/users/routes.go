package users

import (
	"github.com/RiCliZz/shop/responses"
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
	store          types.UserStore
	storeAddresses types.AddressStore
}

func NewHandler(store types.UserStore, storeAddresses types.AddressStore) *Handler {
	return &Handler{store: store, storeAddresses: storeAddresses}
}

func (h *Handler) RegisterRoutesUser(router *mux.Router) {
	router.HandleFunc("/register", h.registerHandler).Methods("POST")
	router.HandleFunc("/login", h.loginHandler).Methods("POST")
	router.HandleFunc("/confirm", h.confirmEmailHandler).Methods("GET")
	router.HandleFunc("/profile", auth.WithJWTAuth(h.profileHandler, h.store)).Methods("GET")
	router.HandleFunc("/profile", auth.WithJWTAuth(h.updateProfileHandler, h.store)).Methods("PATCH")

	//only for admin
	router.HandleFunc("/user/{id}", auth.WithJWTAdminAuth(h.handleGetUser, h.store)).Methods("GET")
}

// @Summary Update Profile
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body types.UserUpdatePayload true "User data"
// @Success 200 {object} responses.SuccessResponse "Success update"
// @Router /profile [PATCH]
func (h *Handler) updateProfileHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Error: "invalid user id",
		})
		return
	}
	u, err := h.store.GetUserByIDForProfile(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "user does not exist",
			Details: err.Error(),
		})
		return
	}
	var user types.UserUpdatePayload
	if err = utils.ParseJSON(r, &user); err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Error:   "invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err = utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: errors.Error(),
			Error:   "validation error",
		})
		return
	}

	if user.FirstName != nil {
		u.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		u.LastName = *user.LastName
	}
	if user.Password != nil {
		hashPass, err := auth.HashPass(*user.Password)
		if err != nil {
			utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
				Error:   "failed to hash password",
				Details: err.Error(),
			})
			return
		}
		u.Password = hashPass
	}
	if err = h.store.UpdateUserProfile(id, u); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   "failed to update user profile",
			Details: err.Error(),
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, responses.SuccessResponse{
		Success: true,
		Data:    u,
	})
}

func (h *Handler) confirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	err := h.store.CheckToken(r.URL.Query().Get("token"))
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error_confirm_email",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, "SUCCESS")
}

// @Summary Get MY Profile
// @Description Profile for User
// @Tags Users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} types.UserProfile "Get your profile"
// @Router /profile [GET]
func (h *Handler) profileHandler(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: "Not authorized",
		})
		return
	}
	u, err := h.store.GetUserByIDForProfile(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   err.Error(),
			Details: "not found user with id",
		})
		return
	}
	address, err := h.storeAddresses.GetAddresses(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   err.Error(),
			Details: "not found user with id",
		})
	}
	userProfile := types.UserProfile{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
		Address:   address,
	}
	utils.WriteJSON(w, http.StatusOK, userProfile)
}

// @Summary log in
// @Description Log in user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body types.UserLoginPayload true "User data"
// @Success 200 {object} responses.JWTResponse "Success login"
// @Router /login [post]
func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	//Приведение к JSON
	var user types.UserLoginPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "invalid_request",
		})
		return
	}
	//Проверка валидатором
	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: errors.Error(),
			Error:   "validation failed",
		})
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "not found user with this email",
		})
		return
	}
	if !auth.ComparePass(u.Password, user.Password) {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: "password not match",
		})
		return
	}
	if !u.Email_verified {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: "email not verified",
		})
		return
	}
	godotenv.Load()
	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := auth.CreateJWT(secret, u.Id, u.Role)
	utils.WriteJSON(w, http.StatusOK, responses.JWTResponse{token})
}

// @Summary sign in
// @Description Create someone user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body types.UserRegisterPayload true "User data"
// @Success 201 {object} responses.UserRegisterResponse "Success registration, but need confirm Email"
// @Router /register [post]
func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	//Приведение к JSON
	var user types.UserRegisterPayload
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "invalid request body",
		})
		return
	}

	//Проверка валидатором
	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: errors.Error(),
			Error:   "validation error",
		})
		return
	}
	//Проверка на уже существующий email
	_, err := h.store.GetUserByEmail(user.Email)
	if err == nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error: "user with this email already exists",
		})
		return
	}
	//Хэшируем пароль
	hashedPass, err := auth.HashPass(user.Password)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "hashing password failed",
		})
		return
	}
	token, err := utils.GenerateUUID()
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error creating UUID",
		})
		return
	}
	//Если всё файн - регистрируем
	createdUser, err := h.store.CreateAcc(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
		Password:  hashedPass,
	})
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Details: err.Error(),
			Error:   "error creating user",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdUser)

}

// @Summary Get
// @Description GET User ONLY FOR ADMIN
// @Tags Admin
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "User data"
// @Success 200 {object} types.User "User with this ID: "
// @Router /user/{id} [GET]
func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Details: "id not found",
		})
		return
	}
	id, err := strconv.Atoi(str)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, responses.ErrorResponse{
			Error:   err.Error(),
			Details: "invalid id",
		})
		return
	}
	user, err := h.store.GetUserByID(id)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, responses.ErrorResponse{
			Error:   err.Error(),
			Details: "not found this id",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, user)
}
