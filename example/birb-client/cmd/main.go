package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mattsp1290/october-talks-2025/example/birb-client/internal/config"
	"github.com/mattsp1290/october-talks-2025/example/birb-client/internal/writer"
)

func main() {
	// Load configuration from environment
	cfg := config.Load()

	log.Printf("Birb-Client starting...")
	log.Printf("Configuration:")
	log.Printf("  - Birb Nest URL: %s", cfg.BirbNestURL)
	log.Printf("  - Write Interval: %v", cfg.WriteInterval)
	log.Printf("  - Log Level: %s", cfg.LogLevel)

	// Test connectivity to birb-nest
	healthURL := fmt.Sprintf("%s/health", cfg.BirbNestURL)
	resp, err := http.Get(healthURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to birb-nest: %v", err)
		log.Println("Continuing anyway - service may not be available yet")
	} else {
		resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			log.Println("Successfully connected to birb-nest")
		} else {
			log.Printf("Warning: birb-nest health check returned status %d", resp.StatusCode)
		}
	}

	// Create writer
	w := writer.New(cfg.BirbNestURL, cfg.WriteInterval)

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle SIGINT and SIGTERM
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start writer in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- w.Start(ctx)
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
		log.Println("Initiating graceful shutdown...")
		cancel()
	case err := <-errChan:
		if err != nil && err != context.Canceled {
			log.Printf("Writer error: %v", err)
		}
	}

	// Close writer
	if err := w.Close(); err != nil {
		log.Printf("Error closing writer: %v", err)
	}

	log.Println("Birb-Client shutdown complete")
}