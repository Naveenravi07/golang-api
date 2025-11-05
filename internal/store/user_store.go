package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextpassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextpassword), 12)
	if err != nil {
		return err
	}
	p.hash = hash
	p.plaintext = &plainTextpassword
	return nil
}

func (p *password) Matches(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err //internal server error
		}
	}
	return true, nil
}

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `INSERT INTO users (username,email,password_hash,bio) VALUES ($1,$2,$3,$4) RETURNING id`
	err := pg.db.QueryRow(query, user.Username, user.Email, user.PasswordHash.plaintext, user.Bio).Scan(&user.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(usermame string) (*User, error) {
	user := &User{}
	query := `SELECT id,username,email,password_hash,bio,createdAT,updatedAt from users where username=$1`
	err := pg.db.QueryRow(query, usermame).Scan(
		&user.Id, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Bio, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil

}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `UPDATE users SET username=$1,email=$2,bio=$3,updatedAt=CURRENT_TIMESTAMP WHERE id=$4 RETURNING updatedAt`
	result, err := pg.db.Exec(query, user.Username, user.Email, user.Bio, user.Id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No user found and updated")
	}
	return nil
}
