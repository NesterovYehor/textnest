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
	dbConn, cleanup := PreparePostgres(ctx, t)
	defer cleanup()
	model := models.NewUserModel(dbConn)
	srv := services.NewUserService(model)

	user := &models.User{
		Name:  "test",
		Email: "test@email",
	}

	userID, err := srv.CreateNewUser(user.Name, user.Email, "test-password")
	assert.NoError(t, err, "Error creating new user")

	exist, err := srv.UserExists(userID)
	assert.NoError(t, err, "Error checking if user exists")
	assert.True(t, exist, "User should exist after creation")

	err = srv.ActivateUser(userID)
	assert.NoError(t, err, "expected no error while activating user")

	var activated bool
	err = dbConn.QueryRow("SELECT activated FROM users WHERE id = $1", userID).Scan(&activated)
	assert.NoError(t, err, "expected no error while querying the user")
	assert.True(t, activated, "user should be activated after calling ActivateUser")
}
