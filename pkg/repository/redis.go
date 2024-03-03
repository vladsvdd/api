package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(address, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisClient{client: client}, nil
}

func (rc *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	value, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("%v key does not exist", key)
		}
		return "", fmt.Errorf("failed to get %v from Redis: %v", key, err)
	}
	return value, nil
}

func (rc *RedisClient) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	err := rc.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set qrator_jsid in Redis: %v", err)
	}
	return nil
}
