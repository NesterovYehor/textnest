package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

// DownloadPasteRequest represents the request body to download a paste.
// @Summary Download a paste
// @Description Retrieves a paste's metadata and download URL by its key
// @Tags paste
// @Accept json
// @Produce json
// @Param input body DownloadPasteRequest true "Paste Key"
// @Success 200 {object} map[string]interface{} "Successful response"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /paste/download [post]
type DownloadPasteRequest struct {
	Key string `json:"key"`
}

// DownloadPaste downloads a paste by its key.
// @Summary Download a paste
// @Description Retrieves a paste's metadata and download URL by its key
// @Tags paste
// @Accept json
// @Produce json
// @Param input body DownloadPasteRequest true "Paste Key"
// @Success 200 {object} map[string]interface{} "Successful response"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /paste/download [post]
func DownloadPaste(cfg *config.Config, app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input DownloadPasteRequest

		// Read JSON from the request body
		if err := helpers.ReadJSON(w, r, &input); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("failed to read JSON input: %w", err), nil)
			errors.BadRequestResponse(w, http.StatusBadRequest, err)
			return
		}

		// Download paste by key
		downloadResp, err := app.DownloadClient.DownloadByKey(input.Key)
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error downloading paste: %w", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		// Create response
		response := helpers.Envelope{
			"creation_date":   downloadResp.Metadata.CreatedAt.AsTime(),
			"content_url":     downloadResp.DownlaodUrl,
			"title":           downloadResp.Metadata.Title,
			"expiration_date": downloadResp.Metadata.ExpiredDate.AsTime(),
		}

		// Send response
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}

// DownloadAllPastesOfUser retrieves all pastes created by the authenticated user.
// @Summary Get all pastes for a user
// @Description Retrieves a list of pastes for the currently authenticated user
// @Tags paste
// @Accept json
// @Produce json
// @Param limit query int false "Limit the number of results (default: 10)"
// @Param offset query int false "Offset for pagination (default: 0)"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Successful response with pastes"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /pastes [get]
func DownloadAllPastesOfUser(cfg *config.Config, app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId, ok := r.Context().Value("user_id").(string)
		if !ok {
			errors.NoTokenProvided(w)
			app.Logger.PrintError(ctx, fmt.Errorf("Failed to get user id"), nil)
			return
		}

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if limit == "" {
			limit = "10"
		}
		if offset == "" {
			offset = "0"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			http.Error(w, "Invalid limit value", http.StatusBadRequest)
			return
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "Invalid offset value", http.StatusBadRequest)
			return
		}

		// Fetch pastes for user
		pastes, err := app.DownloadClient.DownloadByUserId(userId, int32(limitInt), int32(offsetInt))
		if err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Failed to fetch pastes: %v", err), nil)
			errors.ServerErrorResponse(w, err)
			return
		}

		// Create response
		response := helpers.Envelope{
			"pastes": pastes.Objects,
		}

		// Send response
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}
