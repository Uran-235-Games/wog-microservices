package db

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisDB struct {
	client *redis.Client
}

func (r *RedisDB) Connect(addr string) error {
    r.client = redis.NewClient(&redis.Options{
        Addr:	  addr,
        Password: "", // No password set
        DB:		  0,  // Use default DB
        Protocol: 2,  // Connection protocol
    })
	return nil
}

func (r *RedisDB) Disconnect() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}