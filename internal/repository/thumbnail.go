package repository

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"thumbnail/internal/storage/cache"
)

type Repository struct {
	Storage *cache.Storage
}

type ThumbnailRepository interface {
	SetThumbnail(ctx context.Context, url string, thumbnail []byte) error
	GetThumbnail(ctx context.Context, url string) ([]byte, error)
}

func NewThumbnailRepository(storage *cache.Storage) (*Repository, error) {
	slog.Info("Repository layer successfully create!")
	return &Repository{Storage: storage}, nil
}

func (r *Repository) GetThumbnail(ctx context.Context, url string) ([]byte, error) {
	if url == "" {
		return nil, errors.New("url can't be empty")
	}

	data, err := r.Storage.Rdb.Get(ctx, url).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("thumbnail not found or expired")
	}

	return []byte(data), nil
}

func (r *Repository) SetThumbnail(ctx context.Context, url string, thumbnail []byte) error {
	if url == "" {
		return errors.New("url can't be empty")
	}

	durationStr := os.Getenv("TTL")

	ttl, err := time.ParseDuration(durationStr)
	if err != nil {
		return err
	}

	err = r.Storage.Rdb.Set(ctx, url, thumbnail, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
