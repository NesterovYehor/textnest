package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/validation"
)

// UploadPaste handles uploading a paste with expiration and content.
func UploadPaste(w http.ResponseWriter, r *http.Request, cfg *config.Config, ctx context.Context, app *app.AppContext) {
	// Parse input JSON body
	var input validation.PasteInput
	if err := helpers.ReadJSON(w, r, &input); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("failed to read JSON input: %w", err), nil)
		errors.BadRequestResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format"))
		return
	}

	// Validate input
	if err := validation.ValidatePasteInput(&input); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("validation error: %w", err), nil)
		errors.BadRequestResponse(w, http.StatusBadRequest, err)
		return
	}

	// Generate a new unique key
	key, err := app.KeyGenClient.GetKey(ctx)
	if err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("error generating new key: %w", err), nil)
		errors.ServerErrorResponse(w, fmt.Errorf("internal error while generating key"))
		return
	}

	// Upload the paste
	uploadRes, err := app.UploadClient.Upload(key, input.ExpirationDate, []byte(input.Content))
	if err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("error uploading paste: %w", err), nil)
		errors.ServerErrorResponse(w, fmt.Errorf("internal error while uploading paste"))
		return
	}
	if uploadRes == "" {
		app.Logger.PrintError(ctx, fmt.Errorf("empty response from upload service"), nil)
		errors.ServerErrorResponse(w, fmt.Errorf("internal error: empty response from upload service"))
		return
	}

	// Send success response
	response := helpers.Envelope{"message": uploadRes, "key": key}
	if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
		errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
	}
}
