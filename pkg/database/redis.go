package database

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	addr := "localhost:6379"

	if envAddr := os.Getenv("REDIS_ADDR"); envAddr != "" {
		addr = envAddr
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		// Password: "", // no password set
		// DB:       0,  // use default DB
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Could not connect to Redis at %s: %v", addr, err)
	}

	log.Printf("Connected to Redis at %s", addr)
	return rdb
}
