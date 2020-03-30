package users

import (
	"database/sql"

	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64
	Username string
}

var (
	ErrNoRecord           = errors.New("No matching record found")
	ErrInvalidCredentials = errors.New("Invalid credentials")
)

//go:generate sh -c "mockgen -source=users.go -package users -destination users_mock.go UserService"
type UserService interface {
	Insert(tx *sqlx.Tx, username, password string) (*User, error)
	GetByID(tx *sqlx.Tx, userID int64) (*User, error)
	Authenticate(tx *sqlx.Tx, username, password string) (int64, error)
}

// Implements UserService
type UserManager struct {
}

func (um *UserManager) Insert(tx *sqlx.Tx, username, password string) (*User, error) {
	const q = `
      INSERT INTO users(username, hashed_password)
	  VALUES (lower($1), $2)
	  RETURNING id
    `
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	var rowId int64
	err = tx.QueryRowx(q, username, hashedPassword).Scan(&rowId)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       rowId,
		Username: username,
	}, nil
}

func (um *UserManager) GetByID(tx *sqlx.Tx, userID int64) (*User, error) {
	const q = `
      SELECT id, username
      FROM users
      WHERE id=$1
    `
	var user User
	err := tx.Get(&user, q, userID)
	if err != nil {
		return nil, ErrNoRecord
	}
	return &user, nil
}

func (um *UserManager) Authenticate(tx *sqlx.Tx, username, password string) (int64, error) {
	const q = `
		SELECT id, hashed_password
        FROM users
        WHERE username=lower($1)
	`
	var (
		id             int64
		hashedPassword []byte
	)

	row := tx.QueryRowx(q, username)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, ErrNoRecord
		}
		return -1, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return -1, ErrInvalidCredentials
	} else if err != nil {
		return -1, err
	}
	return id, nil
}
