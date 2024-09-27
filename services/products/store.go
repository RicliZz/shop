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

func (s *Store) GetProductByID(id int) (*types.CreateProductPayload, error) {
	var p types.CreateProductPayload
	err := s.db.QueryRow("SELECT name, description, price FROM products WHERE id = $1", id).Scan(
		&p.Name,
		&p.Description,
		&p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Store) GetProducts() ([]*types.ShortProducts, error) {
	rows, err := s.db.Query("SELECT name, price FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]*types.ShortProducts, 0)
	for rows.Next() {
		product, err := scanProducts(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (s *Store) CreateProduct(p *types.CreateProductPayload) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, price) VALUES ($1, $2, $3)", p.Name, p.Description, p.Price)
	if err != nil {
		return err
	}
	return nil
}

func scanProducts(rows *sql.Rows) (*types.ShortProducts, error) {
	products := new(types.ShortProducts)
	err := rows.Scan(
		&products.Name,
		&products.Price)
	if err != nil {
		return nil, err
	}
	return products, nil
}
