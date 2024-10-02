package orders

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(store *sql.DB) *Store {
	return &Store{db: store}
}

func (s *Store) CreateNewOrder(userID int, total float64) (int, error) {
	var orderID int
	query := `INSERT INTO orders (user_id, total, status, created_at) 
			  VALUES ($1, $2, 'pending', NOW()) RETURNING id`

	err := s.db.QueryRow(query, userID, total).Scan(&orderID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}
