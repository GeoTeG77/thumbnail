package service

import (
	"context"
	"log/slog"
	
	"thumbnail/internal/models"
	"thumbnail/internal/repository"
)

type ThumbnailService struct {
	storage *repository.Repository
}

func NewThumbnailService(repo *repository.Repository) (*ThumbnailService, error) {
	slog.Info("Service layer successfully create!")
	return &ThumbnailService{
		storage: repo,
	}, nil
}

func (s *ThumbnailService) GetThumbnail(ctx context.Context, url string) (*models.ThumbnailResponse, error) {
	thumbnail, err := s.storage.GetThumbnail(ctx, url)
	if err != nil {
		return nil, err
	}

	response := &models.ThumbnailResponse{
		URL:     url,
		Preview: thumbnail,
	}

	return response, nil
}

func (s *ThumbnailService) SetThumbnail(ctx context.Context, url string, thumbnail []byte) error {
	if err := s.storage.SetThumbnail(ctx, url, thumbnail); err != nil {
		return err
	}

	return nil
}
