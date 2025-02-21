package handler

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

// ExpirePasteHandler godoc
// @Summary Expire a single paste by its key
// @Description Expire a specific paste by providing the paste key. Requires user authentication.
// @Tags paste
// @Accept json
// @Produce json
// @Param key path string true "Paste Key"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Paste expired successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
func ExpirePasteHandler(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, ok := ctx.Value("user_id").(string)
		key := r.PathValue("key")

		if !ok {
			app.Logger.PrintError(ctx, fmt.Errorf("Authorization failed: user_id missing"), nil)
			errors.NoTokenProvided(w)
			return
		}

		if key == "" {
			app.Logger.PrintError(ctx, fmt.Errorf("Missing paste key"), nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, fmt.Errorf("Paste key is required"))
			return
		}

		res, err := app.UploadClient.ExpirePaste(ctx, key, userId)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Expiring paste failed: %v", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		response := helpers.Envelope{"message": res}
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("Internal error while sending response"))
		}
	}
}

// ExpireAllUserPastesHandler godoc
// @Summary Expire all pastes of a user
// @Description Expire all pastes associated with the authenticated user. Requires user authentication.
// @Tags paste
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Paste expired successfully"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /expire_all [delete]
func ExpireAllUserPastesHandler(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, ok := ctx.Value("user_id").(string)

		if !ok {
			app.Logger.PrintError(ctx, fmt.Errorf("Authorization failed: user_id missing"), nil)
			errors.NoTokenProvided(w)
			return
		}

		res, err := app.UploadClient.ExpireAllUserPastes(ctx, userId)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Expiring all pastes failed: %v", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		response := helpers.Envelope{"message": res}
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("Internal error while sending response"))
		}
	}
}
