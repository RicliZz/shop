package cart

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

func (s *Store) AddToCart(id_u int, id_p int, quantity int) error {
	_, err := s.db.Exec("INSERT INTO cart (user_id, product_id, quantity) VALUES ($1, $2, $3)", id_u, id_p, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CheckCart(id_u int) ([]types.CartItem, error) {

	rows, err := s.db.Query("SELECT p.name, p.price, c.quantity FROM cart c JOIN products p on c.product_id = p.id WHERE user_id=$1", id_u)
	if err != nil {
		return nil, err
	}
	var products []types.CartItem
	for rows.Next() {
		var product types.CartItem
		if err = rows.Scan(&product.Name, &product.Price, &product.Quantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) DeleteCart(id int) error {
	_, err := s.db.Exec("DELETE FROM cart WHERE user_id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
