package router

import (
	"log"
	"net/http"

	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/handlers"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/storage"
)

func Router(cfg *config.Config) *http.ServeMux {
	s3, err := storage.NewS3Storage(cfg.Storage.Bucket, cfg.Storage.Region)
	if err != nil {
		log.Fatalf("Error creating router: %v", err)
	}

	router := http.NewServeMux()

	router.HandleFunc("POST /v1/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.StorePaste(w, r, s3)
	})

	router.HandleFunc("GET  /v1/fetch/{key}", func(w http.ResponseWriter, r *http.Request) {
		handlers.FetchPaste(w, r, s3)
	})

	router.HandleFunc("DELETE   /v1/fetch/{key}", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeletePaste(w, r, s3)
	})

	router.HandleFunc("/v1/healthcheck", handlers.Healthcheck)

	return router
}
