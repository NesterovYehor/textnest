package integration

import (
	"context"
	"testing"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/NesterovYehor/textnest/services/auth_service/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
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

	dbConn, cleanup := PreparePostgres(ctx, "users", tableSchema, t)
	defer cleanup()
	model := models.NewUserModel(dbConn)
	srv := services.NewUserService(model)

	user := &models.User{
		Name:  "test",
		Email: "test@email",
	}

	userID, err := srv.CreateNewUser(user.Name, user.Email, "test-password")
	assert.NoError(t, err, "Error creating new user")

	exist, err := srv.UserExists(userID.String())
	assert.NoError(t, err, "Error checking if user exists")
	assert.True(t, exist, "User should exist after creation")

	var activated bool
	err = dbConn.QueryRow(ctx, "SELECT activated FROM users WHERE id = $1", userID).Scan(&activated)
	assert.NoError(t, err, "expected no error while querying the user")
}
