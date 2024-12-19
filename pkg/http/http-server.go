package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// NewConfig initializes default server settings
func NewConfig(port string) *Config {
	return &Config{
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func RunServer(ctx context.Context, cfg *Config, handler http.Handler) error {
	srv := http.Server{
		Addr:         cfg.Port,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		<-context.TODO().Done()
		fmt.Println("shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			fmt.Printf("server shutdown error: %v", err)
		}
	}()

	fmt.Printf("server listening on port %s\n", cfg.Port)
	return srv.ListenAndServe()
}
