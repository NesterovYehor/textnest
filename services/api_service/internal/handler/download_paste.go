package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/NesterovYehor/TextNest/pkg/errors"
	"github.com/NesterovYehor/TextNest/pkg/helpers"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	download_service "github.com/NesterovYehor/TextNest/services/api_service/internal/grpc_client/download_service_client"
)

func DownloadPaste(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	key := r.PathValue("key")
	if key == "" {
		errors.IncorrectUrlParams(w, "key")
		return
	}
	fmt.Println("Fetched key: ", key)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	req := download_service.DownloadRequest{
		Key: key,
	}
	res, err := cfg.DownloadService.Download(ctx, &req)
	fmt.Println("Fetched service response: ", res)
	if err != nil {
		errors.ServerErrorResponse(w, err)
		return
	}

	env := helpers.Envelope{
		"Created At": res.CreatedDate.AsTime(),
		"Content":    string(res.Content),
	}

	err = helpers.WriteJSON(w, env, http.StatusOK, nil)
	if err != nil {
		errors.ServerErrorResponse(w, err)
	}
}
