package handler

import (
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	pb "github.com/NesterovYehor/TextNest/services/api_service/api/upload_service"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/validation"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UploadPasteHandler godoc
// @Summary Upload a paste
// @Description Upload a paste with title and expiration date
// @Tags pastes
// @Accept json
// @Produce json
// @Param paste body validation.PasteInput true "Paste Input"
// @Success 200 {object} map[string]interface{} "Upload URL and Key"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /upload [post]

func UploadPasteHandler(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var input validation.PasteInput
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("failed to read JSON input: %w", err), nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, fmt.Errorf("invalid JSON format"))
			return
		}

		app.Logger.PrintInfo(ctx, input.Title, nil)

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
		uploadReq := &pb.UploadPasteRequest{
			UserId:         userID,
			Key:            key,
			Title:          input.Title,
			ExpirationDate: timestamppb.New(input.ExpirationDate),
		}

		uploadURL, err := app.UploadClient.UploadPaste(ctx, uploadReq)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error uploading paste: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while uploading paste"))
			return
		}

		if uploadURL == "" {
			app.Logger.PrintError(ctx, fmt.Errorf("empty response from upload service"), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error: empty response from upload service"))
			return
		}

		response := helpers.Envelope{"upload_url": uploadURL, "key": key}
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}

// UpdatePasteHandler godoc
// @Summary Update a paste
// @Description Update an existing paste based on the key provided
// @Tags pastes
// @Accept json
// @Produce json
// @Param key path string true "Paste Key"
// @Success 200 {object} map[string]interface{} "Updated Paste URL"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /update/{key} [put]
func UpdatePasteHandler(app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		key := r.PathValue("key")

		userId, ok := ctx.Value("user_id").(string)
		if !ok {
			app.Logger.PrintError(ctx, fmt.Errorf("authorization failed: user_id missing"), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("authorization failed"))
			return
		}

		updateRes, err := app.UploadClient.UpdatePaste(
			ctx,
			key,
			userId,
		)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("updating paste failed: %w", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		response := helpers.Envelope{"update_url": updateRes}
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}
