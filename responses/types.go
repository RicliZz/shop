package responses

import (
	"github.com/RiCliZz/shop/types"
	"time"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ProfileResponse struct {
	User    *types.UserProfile `json:"user"`
	Address *types.Address     `json:"address"`
}

type UserRegisterResponse struct {
	Id        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type JWTResponse struct {
	Token string `json:"token"`
}
