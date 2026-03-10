package redis

import (
	"context"
	"fmt"
	"log"
	"redis-learn/config"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewClient(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Printf("Connected to Redis: %s\n", pong)

	return client
}
