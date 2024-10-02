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
	Banned         bool      `json:"banned"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Cart struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ProductID int       `json:"productId"`
	CreatedAt time.Time `json:"createdAt"`
}

type CartItem struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type ShortProducts struct {
	Id    int     `json:"id"`
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
	DeleteAccount(id int) error
	CheckToken(token string) error
	GetUserByIDForProfile(id int) (*UserProfile, error)
	UpdateUserProfile(id int, user *UserProfile) error
	BanUser(id int) error
}

type ProductStore interface {
	GetProductByID(id int) (*CreateProductPayload, error)
	GetProducts() ([]*ShortProducts, error)
	CreateProduct(p *CreateProductPayload) error
	GetProductByName(name string) (*Product, error)
	UpdateProduct(product *Product) error
}

type AddressStore interface {
	CreateNewAddress(id int, store Address) error
	GetAddresses(id int) (*AddressPayload, error)
}

type CartStore interface {
	AddToCart(id_user int, id_product int, quantity int) error
	CheckCart(id int) ([]CartItem, error)
	DeleteCart(id int) error
}

type OrderStore interface {
	CreateNewOrder(id int, total float64) (int, error)
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

type Order struct {
	Id        int       `json:"id"`
	User_id   int       `json:"user_id"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

type AddProduct struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type AddToCart struct {
	Id       int `json:"id"`
	Quantity int `json:"quantity"`
}
