package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/handlers"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/storage"
)

func StartServer(cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	s3, err := storage.NewS3Storage(cfg.Storage.Bucket, cfg.Storage.Region)
	if err != nil {
		fmt.Errorf("Failed create new storage: %v", err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.StorePaste(w, r, s3)
	})

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      nil,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to run server: %v\n", err)
		}
	}()

	// To stop the server gracefully when needed, listen for an OS signal here, and call `srv.Shutdown(ctx)`.
}
