package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateNewPaste(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	defer cancel()

	key, err := getNewKey(cfg, ctx)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}

	res := helpers.Envelope{
		"Key": key,
	}

	err = WriteJSON(w, res, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}

func getNewKey(cfg *config.Config, ctx context.Context) (string, error) {
	conn, err := grpc.NewClient(cfg.Grpc.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}

	defer conn.Close()

	client := key_manager.NewKeyManagerServiceClient(conn)

	req := &key_manager.GetKeyRequest{}
	res, err := client.GetKey(ctx, req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return res.Key, err
}

func WriteJSON(w http.ResponseWriter, data any, status int, headers http.Header) error {
	// Marshal the data into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if headers != nil {
		for key, value := range headers {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
