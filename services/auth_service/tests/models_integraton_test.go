package integration

import (
	"context"
	"testing"
	"time"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/stretchr/testify/assert"
)

// Create user and insert
var user = models.User{
	Name:  "test",
	Email: "test@email",
}

func TestInsertUser(t *testing.T) {
	// Set up container and DB connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tableSchema := `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
            created_at timestamp(0) with time zone NOT NULL DEFAULT NOW (),
            name text NOT NULL,
            email citext UNIQUE NOT NULL,
            password_hash bytea NOT NULL,
            activated bool NOT NULL DEFAULT false
        );
    `

	// Prepare the Postgres container and connection
	conn, cleanup := PreparePostgres(ctx, "users", tableSchema, t)
	defer cleanup()
	userModel := models.NewUserModel(conn)
	err := user.Password.Set("TestPassword")
	assert.NoError(t, err)

	// Insert first user
	userID, err := userModel.Insert(&user)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)

	// Insert user with the same email to test duplicate
	duplicateUser := &models.User{
		Name:     "duplicate",
		Email:    "test@email",
		Password: user.Password,
	}
	_, err = userModel.Insert(duplicateUser)
	assert.Error(t, err)
	assert.Equal(t, models.ErrDuplicateEmail, err)
}

func TestInsertToken(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tableSchema := `
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
            created_at timestamp(0) with time zone NOT NULL DEFAULT NOW (),
            name text NOT NULL,
            email citext UNIQUE NOT NULL,
            password_hash bytea NOT NULL,
            activated bool NOT NULL DEFAULT false
        );
        CREATE TABLE IF NOT EXISTS tokens (
            hash TEXT PRIMARY KEY,
            user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL
        );`
	conn, cleanup := PreparePostgres(ctx, "tokens", tableSchema, t)
	defer cleanup()

	token := models.Token{
		Hash:   "test-hash",
		Expiry: time.Now().Add(time.Minute),
	}
	err := user.Password.Set("Test-password")
	if err != nil {
		t.Fatal(err)
	}

	tokensModel := models.NewTokenModel(conn)
	usersModel := models.NewUserModel(conn)
	userID, err := usersModel.Insert(&user)
	if err != nil {
		t.Fatal(err)
	}
	token.UserID = *userID
	assert.NoError(t, tokensModel.Insert(&token))
}
