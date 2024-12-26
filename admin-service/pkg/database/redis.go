package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"cybermind/admin-service/configs"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis initializes Redis connection
func InitRedis(config *configs.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := RDB.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	log.Println("Redis connected successfully")
	return nil
}

// GetKey retrieves a key's value
func GetKey(ctx context.Context, key string) (string, error) {
	if RDB == nil {
		return "", fmt.Errorf("redis client not initialized")
	}
	return RDB.Get(ctx, key).Result()
}

// SetKey sets a key's value
func SetKey(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	if RDB == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return RDB.Set(ctx, key, value, ttl).Err()
}

// DelKey deletes a key
func DelKey(ctx context.Context, key string) error {
	if RDB == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return RDB.Del(ctx, key).Err()
}

// SetKeyNX sets a key's value if it does not exist
func SetKeyNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	if RDB == nil {
		return false, fmt.Errorf("redis client not initialized")
	}
	return RDB.SetNX(ctx, key, value, ttl).Result()
} 