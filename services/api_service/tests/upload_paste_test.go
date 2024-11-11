package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client/key_manager_client"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client/upload_service_client"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUploadPaste(t *testing.T) {
	// Initialize real gRPC clients
	kgsConn, err := grpc.NewClient("localhost:5555", grpc.WithTransportCredentials(insecure.NewCredentials())) // gRPC address for KeyManager service
	if err != nil {
		t.Fatalf("Failed to connect to Key Manager service: %v", err)
	}
	defer kgsConn.Close()
	keyManagerClient := key_manager.NewKeyManagerServiceClient(kgsConn)

	uploadConn, err := grpc.NewClient("localhost:3489", grpc.WithTransportCredentials(insecure.NewCredentials())) // gRPC address for Upload service
	if err != nil {
		t.Fatalf("Failed to connect to Upload service: %v", err)
	}
	defer uploadConn.Close()
	uploadClient := upload_service.NewUploadServiceClient(uploadConn)

	// Initialize configuration with real clients
	cfg := &config.Config{
		KeyManager:    keyManagerClient,
		UploadService: uploadClient,
	}

	// Prepare test input
	input := map[string]interface{}{
		"expiration_date": time.Now().Add(time.Hour),
		"content":         "sample content",
	}
	jsonData, _ := json.Marshal(input)

	// Create HTTP request and response recorder
	req := httptest.NewRequest("POST", "/upload", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Call the handler function with real services
	handler.UploadPaste(rr, req, cfg, context.Background())

	// Get the response
	res := rr.Result()
	defer res.Body.Close()

	// Assert the response status and message
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	var response map[string]string
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expectedMessage := "Uploaded new paste successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %q; got %q", expectedMessage, response["message"])
	}
}
