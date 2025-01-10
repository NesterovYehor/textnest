package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/NesterovYehor/TextNest/services/cleanup_service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	app, cleanUp, err := app.NewApp(ctx)
	defer cleanUp()

	// Check if the initialization encountered an error
	if err != nil {
		fmt.Printf("App initialization failed: %v\n", err)
		os.Exit(1)
	}
	app.Logger.PrintInfo(ctx, "Starting Scheduler", nil)
	app.Scheduler.Start(ctx, app.Config.ExpirationInterval)
	// Run the Kafka consumer
	go func() {
		if err := app.RunKafkaConsumer(app.Config, ctx); err != nil {
			app.Logger.PrintError(ctx, fmt.Errorf("Error running Kafka consumer: %v", err), nil)
		}
	}()

	// Graceful shutdown setup
	// Listen for system signals (e.g., SIGINT or SIGTERM) for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a termination signal
	<-signalChan

	// Perform any shutdown tasks or logging here if needed
	fmt.Println("Shutting down gracefully...")
}
