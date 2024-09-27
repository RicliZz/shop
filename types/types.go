package types

import (
	"github.com/google/uuid"
	"time"
)

type UserUpdatePayload struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Password  *string `json:"password" validate:"required,min=8,max=32"`
}

type UserRegisterPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=32"`
}

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type User struct {
	Id             int       `json:"id"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	Email          string    `json:"email"`
	Email_verified bool      `json:"email_verified"`
	Token          uuid.UUID `json:"token"`
	Password       string    `json:"password"`
	CreatedAt      time.Time `json:"createdAt"`
	Role           string    `json:"role"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ShortProducts struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required"`
}
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateAcc(User) (*UserRegisterPayload, error)
	CheckToken(token string) error
	GetUserByIDForProfile(id int) (*UserProfile, error)
	UpdateUserProfile(id int, user *UserProfile) error
}

type ProductStore interface {
	GetProductByID(id int) (*CreateProductPayload, error)
	GetProducts() ([]*ShortProducts, error)
	CreateProduct(p *CreateProductPayload) error
}

type AddressStore interface {
	CreateNewAddress(id int, store Address) error
	GetAddresses(id int) (*AddressPayload, error)
}

type Address struct {
	ID        int    `json:"id"`
	User_ID   int    `json:"user_id"`
	City      string `json:"city"`
	Street    string `json:"street"`
	House     int    `json:"house"`
	Apartment int    `json:"apartment"`
}

type AddressPayload struct {
	City      string `json:"city"`
	Street    string `json:"street"`
	House     int    `json:"house"`
	Apartment int    `json:"apartment"`
}

type UserProfile struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
	Address   *AddressPayload
}
