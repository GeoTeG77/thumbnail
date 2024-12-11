package cache

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Rdb *redis.Client
}

func Init() (*Storage, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")

	if redisHost == "" {
		return nil, fmt.Errorf("REDIS_HOST IS NOT SET")
	}

	if redisPort == "" {
		redisPort = os.Getenv("DEFAULT_REDIS_PORT")
	}

	if redisDB == "" {
		redisDB = os.Getenv("DEFAULT_REDIS_DB")
	}

	numDB, err := strconv.Atoi(redisDB)
	if err != nil {
		return nil, err
	}

	Addr := redisHost+":"+redisPort

	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: redisPassword,
		DB:       numDB,
	})

	storage := &Storage{
		Rdb: client,
	}
	ctx := context.Background()
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	slog.Info("Redis connection successfully!")
	slog.Info("Redis init successfully!")
	return storage, nil
}
