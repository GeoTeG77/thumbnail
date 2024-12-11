package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"thumbnail/internal/config"
	"thumbnail/internal/http/handlers"
	"thumbnail/internal/http/router"
	"thumbnail/internal/proto/server"
	"thumbnail/internal/repository"
	"thumbnail/internal/service"
	"thumbnail/internal/storage/cache"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := config.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := cache.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	thumbnailRepo, err := repository.NewThumbnailRepository(db)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail repository: %w", err)
	}

	thumbnailService, err := service.NewThumbnailService(thumbnailRepo)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail service: %w", err)
	}

	handler := handlers.NewThumbnailHandler(thumbnailService)
	router := router.NewRouter(handler)
	errCh := make(chan error, 1)

	go func() {
		if err := server.Run(os.Getenv("GRPC_PORT"), handler); err != nil {
			errCh <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()

	go func() {
		if err := http.ListenAndServe(os.Getenv("HTTP_PORT"), router); err != nil {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	return <-errCh
}
