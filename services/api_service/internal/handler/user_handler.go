package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

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
			app.Logger.PrintInfo(ctx, fmt.Sprintf("response of Autherization: %v", ress), nil)
			app.Logger.PrintError(ctx, err, nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
		}
		response := helpers.Envelope{
			"access_token":  ress.AccessToken,
			"refresh_token": ress.RefreshToken,
		}

		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			errors.ServerErrorResponse(w, err)
		}
	}
}
