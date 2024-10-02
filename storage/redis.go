package storage

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rotisserie/eris"
	"github.com/voxtmault/panacea-shared-lib/config"
)

var redisClient *redis.Client

func validateRedisConfig(cfg *config.RedisConfig) error {
	if cfg.RedisHost == "" {
		return eris.New("redis host is empty")
	}
	if cfg.RedisPort == "" {
		return eris.New("redis port is empty")
	}
	if cfg.RedisPassword == "" {
		return eris.New("redis password is empty")
	}

	return nil
}

func InitRedis(config *config.RedisConfig) error {

	if err := validateRedisConfig(config); err != nil {
		return eris.Wrap(err, "invalid redis configuration")
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       int(config.RedisDBNum),
	})

	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		return eris.Wrap(err, "Init Redis")
	}

	slog.Info("Successfully opened redis connection")
	return nil
}

func CloseRedis() error {
	if err := redisClient.Close(); err != nil {
		return eris.Wrap(err, "Closing redis connection")
	}

	return nil
}

func GetRedisCon() *redis.Client {
	return redisClient
}

// SaveToRedis can also be use to update values in redis
func SaveToRedis(ctx context.Context, key string, value interface{}) error {
	cfg := config.GetConfig()

	if err := redisClient.Set(ctx, key, value, time.Minute*time.Duration(cfg.RedisConfig.RedisExpiration)).Err(); err != nil {
		return eris.Wrap(err, "saving data to redis cache")
	}

	return nil
}

func GetFromRedis(ctx context.Context, key string) (any, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found")
	} else if err != nil {
		return nil, err
	}

	return value, nil
}
