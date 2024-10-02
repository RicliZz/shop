package users

import (
	"database/sql"
	"fmt"
	"github.com/RiCliZz/shop/types"
	"github.com/RiCliZz/shop/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateAcc(user types.User) (*types.UserRegisterPayload, error) {
	createdUser := new(types.UserRegisterPayload)
	err := s.db.QueryRow(`
		INSERT INTO users (firstname, lastname, email, token, password) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING firstname, lastname, email, password`,
		user.FirstName, user.LastName, user.Email, user.Token, user.Password).Scan(
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Email,
		&createdUser.Password)

	if err != nil {
		return nil, err
	}
	err = utils.EmailSend(user.Email, user.Token)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

func (s *Store) DeleteAccount(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec(`DELETE FROM orders WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM cart WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM address WHERE user_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Store) CheckToken(token string) error {
	u := new(types.User)
	err := s.db.QueryRow("SELECT id FROM users WHERE token=$1", token).Scan(&u.Id)
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	_, err = s.db.Exec("UPDATE users SET email_verified=TRUE WHERE id=$1", u.Id)
	return err
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	u := new(types.User)
	for rows.Next() {
		u, err = scanRows(rows)
		if err != nil {
			return nil, err
		}
	}
	if u.Id == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) GetUserByIDForProfile(id int) (*types.UserProfile, error) {
	u := new(types.UserProfile)
	err := s.db.QueryRow("SELECT firstname, lastname, email, password FROM users WHERE id = $1", id).Scan(
		&u.FirstName, &u.LastName, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	u := new(types.User)
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&u.Id, &u.FirstName, &u.LastName,
		&u.Email, &u.Email_verified, &u.Token, &u.Password, &u.CreatedAt, &u.Role, &u.Banned)
	if err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) UpdateUserProfile(id int, u *types.UserProfile) error {
	_, err := s.db.Exec(`UPDATE users SET firstname=$1, lastname=$2, password=$3 WHERE id=$4`, u.FirstName, u.LastName, u.Password, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) BanUser(id int) error {
	query := `UPDATE users SET banned=TRUE WHERE id=$1`
	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func scanRows(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Email_verified,
		&user.Token,
		&user.Password,
		&user.CreatedAt,
		&user.Role,
		&user.Banned,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
