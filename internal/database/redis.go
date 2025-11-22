package database

import (
	"context"
	"fmt"
	"site-constructor/configs"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *configs.RedisConfig) (*redis.Client, error) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DBNumber,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed connect to Redis: %w", err)
	}

	return rdb, nil
}
