package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"thumbnail/internal/config"
	"thumbnail/internal/http/handlers"
	"thumbnail/internal/http/router"
	"thumbnail/internal/proto/server"
	"thumbnail/internal/repository"
	"thumbnail/internal/service"
	cache "thumbnail/internal/storage/cache"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := config.LoadConfig()
	if err != nil {
		slog.Debug("Failed to init config file")
		slog.Error(err.Error())
		os.Exit(1)
	}

	db, err := cache.Init()
	if err != nil {
		slog.Debug("Failed to init DB")
		slog.Error(err.Error())
		os.Exit(1)
	}

	repository, err := repository.NewThumbnailRepository(db)
	if err != nil {
		slog.Debug("Failed to create thumbnail repository")
		slog.Error(err.Error())
		os.Exit(1)
	}

	service, err := service.NewThumbnailService(repository)
	if err != nil {
		slog.Debug("Failed to create thumbnail service")
		slog.Error(err.Error())
		os.Exit(1)
	}

	handler := handlers.NewThumbnailHandler(service)
	router := router.NewRouter(handler)

	go func() {
		if err := server.Run(":50051", handler); err != nil {
			log.Fatalf("Error running gRPC server: %v", err)
		}
	}()

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}