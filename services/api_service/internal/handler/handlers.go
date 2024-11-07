package handler

import (
	"context"
	"net/http"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	key_manager "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client"
	"google.golang.org/grpc"
)

func CreateNewPaste(w http.ResponseWriter, r *http.Request) {
	var input struct {
		content   string `json: "post_text`
		expiredAt string `json: "expired_at"`
		createdAt string `json: "created_at"`
	}

	err := helpers.ReadJSON(w, r, input)
	if err != nil {
		errors.BadRequestResponse(w, http.StatusBadRequest, err)
	}
}

func getNewKey(cfg *config.Config, ctx context.Context) (string, error) {
	conn, err := grpc.NewClient(cfg.Grpc.Addr)
	if err != nil {
		return "", err
	}

	defer conn.Close()

	client := key_manager.NewKeyManagerServiceClient(conn)

	req := &key_manager.GetKeyRequest{}
	res, err := client.GetKey(ctx, req)
	if err != nil {
		return "", err
	}

	return res.Key, err
}
