package handler

import (
	"context"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

// SignUpHandler godoc
// @Summary Sign up a new user
// @Description Sign up a new user by providing their name, email, and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param name body string true "User Name"
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 200 {string} string "User created"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
func SignUpHandler(app *app.AppContext, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name     string `json:"name"`
			Emain    string `json:"email"`
			Password string `json:"password"`
		}

		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}

		_, err := app.AuthClient.SignUp(input.Name, input.Emain, input.Password)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
		}

		if err := helpers.WriteJSON(w, "User created", http.StatusOK, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}

// LogInHandler godoc
// @Summary Log in to the application
// @Description Log in to the application using email and password to receive an access token and refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "User Email"
// @Param password body string true "User Password"
// @Success 200 {object} map[string]interface{} "Tokens and expiration"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /login [post]
func LogInHandler(app *app.AppContext, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Emain    string `json:"email"`
			Password string `json:"password"`
		}

		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}

		ress, err := app.AuthClient.LogIn(input.Emain, input.Password)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
		}
		response := helpers.Envelope{
			"access_token":       ress.AccessToken,
			"refresh_token":      ress.RefreshToken,
			"expires_at":         ress.ExpiresIn.AsTime(),
			"refresh_expires_at": ress.RefreshExpiresAt.AsTime(),
		}

		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}

// RefreshTokens godoc
// @Summary Refresh the access token using the refresh token
// @Description Use the refresh token to get a new access token and refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh Token"
// @Success 200 {object} map[string]interface{} "New Tokens and expiration"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /refresh [post]
func RefreshTokens(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			Refresh string `json:"refresh_token"`
		}
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}
		ress, err := app.AuthClient.RefreshTokens(input.Refresh)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		response := helpers.Envelope{
			"access_token":       ress.AccessToken,
			"refresh_token":      ress.RefreshToken,
			"expires_at":         ress.ExpiresIn.AsTime(),
			"refresh_expires_at": ress.RefreshExpiresAt.AsTime(),
		}

		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}

// ActivateUser godoc
// @Summary Activate a user account
// @Description This endpoint allows a user to activate their account using the provided JWT token.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token path string true "JWT token for user activation"
// @Success 202 {object} map[string]string "User activated successfully"
// @Failure 400 {object} map[string]string "Invalid token"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /activate/{token} [post]
func ActivateUser(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.PathValue("id")
		message, err := app.AuthClient.ActivateUser(token)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.ServerErrorResponse(w, err)
			return
		}
		response := helpers.Envelope{
			"message": message,
		}
		if err := helpers.WriteJSON(w, response, http.StatusAccepted, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}

func SendPasswordResetEmail(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			Email string `json:"email"`
		}
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}
		message, err := app.AuthClient.SendPasswordResetToken(input.Email)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.ServerErrorResponse(w, err)
			return
		}
		response := helpers.Envelope{
			"message": message,
		}
		if err := helpers.WriteJSON(w, response, http.StatusAccepted, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}

func ResetPassword(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.PathValue("token")
		var input struct {
			Password string `json:"password"`
		}
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}
		message, err := app.AuthClient.ResetPassword(input.Password, token)
		if err != nil {
			app.Logger.PrintError(ctx, err, nil)
			errors.ServerErrorResponse(w, err)
			return
		}
		response := helpers.Envelope{
			"message": message,
		}
		if err := helpers.WriteJSON(w, response, http.StatusAccepted, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}
