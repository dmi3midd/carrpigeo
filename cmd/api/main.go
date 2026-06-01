package main

import (
	"carrpigeo/internal/config"
	"carrpigeo/internal/database"
	"carrpigeo/internal/server"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		slog.Info(
			"server forced to shutdown with error",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	server := server.NewServer(cfg, db)

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	slog.Info(
		"server is running",
		slog.String("address", cfg.HTTPServer.Address),
	)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error(
			"failed to run server",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete")
}
