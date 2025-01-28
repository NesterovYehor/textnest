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
	container, conn, err := PreparePostgres(ctx)
	if err != nil {
		t.Fatalf("Failed to prepare Postgres: %v", err)
	}
	// Ensure the container is terminated after the test
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Errorf("Failed to terminate Postgres container: %v", err)
		}
	}()
	defer conn.Close() // Ensure the connection is closed after the test

	// Create user and insert
	user := &models.User{
		Name:  "test",
		Email: "test@email",
	}

	userModel := models.NewUserModel(conn)
	err = user.Password.Set("TestPassword")
	assert.NoError(t, err)

	// Insert first user
	assert.NoError(t, userModel.Insert(user))

	// Insert user with the same email to test duplicate
	duplicateUser := &models.User{
		Name:     "duplicate",
		Email:    "test@email",
		Password: user.Password,
	}
	err = userModel.Insert(duplicateUser)
	assert.Error(t, err)
	assert.Equal(t, models.ErrDuplicateEmail, err)
}
