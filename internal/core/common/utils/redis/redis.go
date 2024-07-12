package redis

import (
	"Eccomerce-website/internal/core/entity"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func InitRedis(serviceLogger *zap.Logger) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		errorResponse := entity.AppConnectionError.Wrap(err, "failed to connect to Redis")
		serviceLogger.Error("connection error",
			zap.String("timestamp", time.Now().Format(time.RFC3339)),
			zap.String("layer", "serviceLayer"),
			zap.String("function", "InitRedis"),
			zap.Error(errorResponse),
			zap.Stack("stacktrace"),
		)
		fmt.Fprintf(os.Stderr, "error while connecting redis: %s", errorResponse.Error())
		os.Exit(1)
	}

	return &RedisClient{
		client: client,
	}
}

func (r *RedisClient) GetRedisClient() *redis.Client {
	return r.client
}
