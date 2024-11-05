package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := httpserver
	log.Println("Starting KGS server...")
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
