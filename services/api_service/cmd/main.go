package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	httpserver "github.com/NesterovYehor/TextNest/pkg/http"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/config"
	"github.com/NesterovYehor/TextNest/services/api_service/internal/handler"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg := config.InitConfig()
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/upload", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		handler.UploadPaste(w, r, cfg, ctx) // Pass the pointer to cfg
	})
	mux.HandleFunc("/v1/download/{key}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received")
		handler.DownloadPaste(w, r, cfg) // Pass the pointer to cfg
	})

	if err := httpserver.RunServer(&ctx, cfg.Http, mux); err != nil {
		log.Panic(err)
	}
}
