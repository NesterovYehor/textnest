package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/storage_service/router"
)

func StartServer(cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done()

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      router.Router(cfg),
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	log.Printf("Starting server, port%s", cfg.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Failed to run server: %v\n", err)
	}

	// To stop the server gracefully when needed, listen for an OS signal here, and call `srv.Shutdown(ctx)`.
}
