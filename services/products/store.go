package products

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

func (s *Store) GetProductByID(id int) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	p := new(types.Product)
	for rows.Next() {
		p, err = scanProducts(rows)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]*types.Product, 0)
	for rows.Next() {
		product, err := scanProducts(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func scanProducts(rows *sql.Rows) (*types.Product, error) {
	products := new(types.Product)
	err := rows.Scan(
		&products.ID,
		&products.Name,
		&products.Description,
		&products.Price)
	if err != nil {
		return nil, err
	}
	return products, nil
}
