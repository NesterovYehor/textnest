package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidUUID    = errors.New("invalid UUID")
	ErrDatabaseError  = errors.New("database error")
	ErrUpdateFailed   = errors.New("update failed")
	ErrInsertFailed   = errors.New("insert failed")
	ErrSelectFailed   = errors.New("select failed")
	AnonymousUser     = &User{}
)

// User represents a user in the system.
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
}

type password struct {
	Plaintext string
	Hash      []byte
}

// IsAnonymous checks if a user is anonymous.
func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}

// UserModel represents the user model.
type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{db: db}
}

func (model *UserModel) Insert(user *User) (string, error) {
	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	args := []any{user.Name, user.Email, user.Password.Hash}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := model.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return "", ErrDuplicateEmail
		}
		return "", fmt.Errorf("%w: %v", ErrInsertFailed, err)
	}

	return user.ID.String(), nil
}

func (model *UserModel) UserExists(userId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := uuid.Parse(userId)
	if err != nil {
		return false, ErrInvalidUUID
	}

	var exists bool
	err = model.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, ErrRecordNotFound
		}
		return false, fmt.Errorf("%w: %v", ErrSelectFailed, err)
	}

	return exists, nil
}

func (model *UserModel) GetByEmail(email string) (*User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated
        FROM users 
        WHERE email = $1
    `

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var user User
	err := model.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrSelectFailed, err)
	}

	return &user, nil
}

func (m *UserModel) ActivateUser(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
        UPDATE users SET activated = true 
        WHERE id = $1 AND activated = false
    `

	id, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidUUID
	}

	res, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m *UserModel) ResetPassword(plainText, tokenHash string) (*uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	var user User
	err := user.Password.Set(plainText)
	if err != nil {
		return nil, err
	}

	query := `
        UPDATE users 
        SET password_hash = $1
        WHERE id = (
            SELECT user_id FROM tokens 
            WHERE hash = $2 AND expiry > NOW()
        )
        RETURNING id
        `

	args := []any{
		user.Password.Hash,
		tokenHash,
	}

	err = m.db.QueryRowContext(ctx, query, args...).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return &user.ID, nil
}

// Set hashes and sets the password.
func (password *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	password.Plaintext = plaintext
	password.Hash = hash

	return nil
}

func (password *password) Match(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(password.Hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
