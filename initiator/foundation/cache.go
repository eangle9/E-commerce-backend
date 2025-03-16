package foundation

import (
	"Eccomerce-website/internal/constant/state"
	"context"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func InitRedis(log Logger, config state.RedisConfig) *redis.Client {
	//url := "redis://user:password@localhost:6379/0?protocol=3"
	opts, err := redis.ParseURL(config.RedisURL)
	if err != nil {
		log.Fatal(context.Background(), "unable to parse redis url",
			zap.Error(err))
		return nil
	}
	client := redis.NewClient(opts)

	// Ping the Redis server to check the connection
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(context.Background(), "unable to connect to Redis server",
			zap.Error(err))
		return nil
	}
	return client
}
