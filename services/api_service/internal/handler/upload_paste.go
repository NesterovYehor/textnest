package handler

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/validation"
)

func UploadPasteHandler(cfg *config.Config, app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input validation.PasteInput
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("failed to read JSON input: %w", err), nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format"))
			return
		}

		if err := validation.ValidatePasteInput(&input); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("validation error: %w", err), nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}

		key, err := app.KeyGenClient.GetKey(ctx)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error generating new key: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while generating key"))
			return
		}

		userID, ok := ctx.Value("user_id").(string)
		if !ok {
			userID = ""
		}

		uploadRes, err := app.UploadClient.UploadPaste(ctx, key, userID, input.ExpirationDate, []byte(input.Content))
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

		response := helpers.Envelope{"message": uploadRes, "key": key}
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}

func Update ff
