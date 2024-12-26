package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/app"
)

func DownloadPaste(w http.ResponseWriter, r *http.Request, cfg *config.Config, ctx context.Context, app *app.AppContext) {
	var input struct {
		Key string `json:"key"`
	}

	if err := helpers.ReadJSON(w, r, &input); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("failed to read JSON input: %w", err), nil)
		errors.BadRequestResponse(w, http.StatusBadRequest, err)
		return
	}

	downloadResp, err := app.DownloadClient.Download(input.Key)
	if err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("error downloading paste: %w", err), nil)
		errors.ServerErrorResponse(w, err)
		return
	}
	response := helpers.Envelope{
		"creation_date":   downloadResp.CreatedDate.AsTime(),
		"content":         string(downloadResp.Content),
		"expiration_date": downloadResp.ExpirationDate.AsTime(),
	}
	if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
		app.Logger.PrintError(ctx, fmt.Errorf("error writing JSON response: %w", err), nil)
		errors.ServerErrorResponse(w, fmt.Errorf("internal error while sending response"))
	}
}
