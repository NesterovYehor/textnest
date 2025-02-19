package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

var responseBufferPool = sync.Pool{
	New: func() interface{} {
		// Allocate a new bytes.Buffer. Adjust initial capacity if needed.
		return new(bytes.Buffer)
	},
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
		buf := responseBufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		defer responseBufferPool.Put(buf)

		// Create response data
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
		for _, paste := range pastes.Objects {
			app.Logger.PrintInfo(ctx, paste.Title, nil)
		}

		app.Logger.PrintInfo(ctx, fmt.Sprintln(pastes), nil)
		// Directly use the Protobuf message inside the response envelope
		response := helpers.Envelope{
			"pastes": pastes.Objects, // Use GetObjects() to get []Metadata directly
		}
		app.Logger.PrintInfo(ctx, fmt.Sprintln(response), nil)

		// Use WriteJSON, which should handle Protobuf struct conversion automatically
		if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
			errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
		}
	}
}
