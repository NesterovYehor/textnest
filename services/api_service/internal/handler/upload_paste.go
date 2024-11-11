package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client/key_manager_client"
	upload_service "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client/upload_service_client"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UploadPaste handles uploading a paste with expiration and content.
func UploadPaste(w http.ResponseWriter, r *http.Request, cfg *config.Config, ctx context.Context) {
	// Parse input JSON body
	var input struct {
		ExpirationDate time.Time `json:"expiration_date"`
		Content        string    `json:"content"`
	}
	if err := helpers.ReadJSON(w, r, &input); err != nil {
		fmt.Println("Failed to read JSON input:", err)
		errors.BadRequestResponse(w, http.StatusBadRequest, err)
		return
	}

	// Generate a new unique key
	KgsResult, err := cfg.KeyManager.GetKey(ctx, &key_manager.GetKeyRequest{})
	if err != nil {
		fmt.Println("Error getting new key:", err)
		errors.ServerErrorResponse(w, err)
		return
	}

	req := &upload_service.UploadRequest{
		Key:            KgsResult.Key,
		ExpirationDate: timestamppb.New(input.ExpirationDate),
		Data:           []byte(input.Content),
	}

	// Execute upload request
	res, err := cfg.UploadService.Upload(ctx, req)
	if err != nil {
		fmt.Println("Error uploading data:", err)
		errors.ServerErrorResponse(w, err)
		return
	}
	if res == nil {
		errors.ServerErrorResponse(w, fmt.Errorf("received nil response from upload service"))
		return
	}

	// Send success response
	response := helpers.Envelope{"message": res.Message, "Key": KgsResult.Key}
	if err := helpers.WriteJSON(w, response, http.StatusOK, nil); err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
