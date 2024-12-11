package repository

import (
	"context"
	"errors"
	"log/slog"
	cache "thumbnail/internal/storage/cache"
	"time"
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
		return nil, errors.New("videoID must not be empty")
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
		return errors.New("videoID must not be empty")
	}

	ttl := 5 * time.Minute
	err := r.Storage.Rdb.Set(ctx, url, thumbnail, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
