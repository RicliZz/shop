package address

import (
	"database/sql"
	"github.com/RiCliZz/shop/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateNewAddress(user_id int, address types.Address) error {
	_, err := s.db.Exec("INSERT INTO address (user_id, city, street, house, apartment) VALUES ($1, $2, $3, $4, $5)",
		user_id, address.City, address.Street, address.House, address.Apartment)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetAddresses(user_id int) (*types.AddressPayload, error) {
	var address types.AddressPayload
	err := s.db.QueryRow("SELECT city, street, house, apartment FROM address WHERE user_id = $1", user_id).Scan(
		&address.City,
		&address.Street,
		&address.House,
		&address.Apartment)
	if err != nil {
		return nil, err
	}
	return &address, nil
}
