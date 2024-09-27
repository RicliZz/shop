package auth

import (
	"context"
	"fmt"
	"github.com/RiCliZz/shop/responses"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.GetTokenFromRequest(r)
		if err != nil {
			log.Println("error getting token from request:", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("error validating token:", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		id := claims["userID"].(string)
		userID, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user: %v", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", u.Id)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
	return f
}

func WithJWTAdminAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := utils.GetTokenFromRequest(r)
		if err != nil {
			log.Println("error getting token from request:", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "Not auth!",
			})
			return
		}
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Println("error validating token:", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		id := claims["userID"].(string)
		role := claims["roleType"].(string)
		if role != "admin" {
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Error:   "Only for admin. Sorry :)",
				Details: "Permission denied",
			})
			return
		}
		userID, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("failed to convert userID to int: %v", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "permission denied",
			})
			return
		}
		u, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user: %v", err)
			utils.ErrorJSON(w, http.StatusForbidden, responses.ErrorResponse{
				Details: err.Error(),
				Error:   "failed to get user",
			})
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", u.Id)
		r = r.WithContext(ctx)
		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
}

func CreateJWT(secret []byte, userID int, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return "", err
	}
	dur, err := strconv.Atoi(os.Getenv("JWT_EXP"))
	if err != nil {
		log.Fatal("Error parsing JWT_EXP")
		return "", err
	}
	expiration := time.Second * time.Duration(int64(dur))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"roleType":  role,
		"expiredAt": time.Now().Add(expiration).Unix(),
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
