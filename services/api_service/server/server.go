package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NesterovYehor/pastebin/tree/main/internal/services/api_service/config"
)

// RunServer starts the HTTP server and handles graceful shutdowns.
func RunServer(cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done() // Mark this function as done when it exits

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Addr),
		Handler:      nil, // Set your handler here
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	// Signal handling goroutine
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit // Wait for a signal
		log.Printf("Caught signal: %+v", s)

		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Attempt to shutdown the server gracefully
		if err := srv.Shutdown(ctx); err != nil {
			shutdownError <- fmt.Errorf("server shutdown failed: %v", err)
		}
		shutdownError <- nil
	}()

	// Start the HTTP server
	go func() {
		log.Printf("Starting server on %s in %s environment...\n", srv.Addr, cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error running server: %s", err)
			shutdownError <- err // Signal shutdown error
		}
	}()

	// Wait for shutdown signal or error
	if err := <-shutdownError; err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	// Wait for all tasks to finish (if any)
	wg.Wait()
	log.Println("Server has shut down gracefully.")
}
