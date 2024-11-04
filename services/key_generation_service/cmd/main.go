package main

import (
	"log"

	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/config"
	"github.com/NesterovYehor/TextNest/tree/main/services/key_generation_service/internal/server"
)

func main() {
	cfg := config.InitConfig()
	srv := server.Server{
		Config: cfg,
	}

	log.Println("Starting KGS server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

