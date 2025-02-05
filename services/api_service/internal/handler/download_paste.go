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

func DownloadHandler(cfg *config.Config, app *app.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app.Logger.PrintInfo(ctx, fmt.Sprintln(ctx), nil)
		userID := ctx.Value("user_id") // Use the correct context key

		if userID == nil {
			app.Logger.PrintInfo(ctx, "No user ID found, calling DownloadPaste", nil)
			DownloadPaste(cfg, app)(w, r)
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			app.Logger.PrintInfo(ctx, "User ID is empty or invalid, calling DownloadPaste", nil)
			DownloadPaste(cfg, app)(w, r)
			return
		}

		app.Logger.PrintInfo(ctx, fmt.Sprintf("User authenticated (ID: %s), calling DownloadAllPastesOfUser", userIDStr), nil)
		DownloadAllPastesOfUser(userIDStr, cfg, app)(w, r)
	})
}

func DownloadPaste(cfg *config.Config, app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var input struct {
			Key string `json:"key"`
		}

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

		// Create response data
		response := helpers.Envelope{
			"creation_date":   downloadResp.CreatedDate.AsTime(),
			"content":         string(downloadResp.Content),
			"expiration_date": downloadResp.ExpirationDate.AsTime(),
		}

		// Send response
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}

func DownloadAllPastesOfUser(userId string, cfg *config.Config, app *app.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		// Set default values if not provided
		if limit == "" {
			limit = "10" // Default limit value
		}
		if offset == "" {
			offset = "0" // Default offset value
		}

		// Convert limit and offset to integers
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
			errors.
				ServerErrorResponse(w, err)
			return
		}

		// Prepare response
		response := helpers.Envelope{
			"pastes": pastes,
		}

		// Send response
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}
