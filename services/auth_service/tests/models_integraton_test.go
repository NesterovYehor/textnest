package integration

import (
	"context"
	"testing"

	"github.com/NesterovYehor/textnest/services/auth_service/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	// Set up container and DB connection
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Prepare the Postgres container and connection
	conn, cleanup := PreparePostgres(ctx, t)
	defer cleanup()

	// Create user and insert
	user := &models.User{
		Name:  "test",
		Email: "test@email",
	}

	userModel := models.NewUserModel(conn)
	err := user.Password.Set("TestPassword")
	assert.NoError(t, err)

	// Insert first user
	userID, err := userModel.Insert(user)
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
	assert.NoError(t, userModel.ActivateUser(userID))
}
