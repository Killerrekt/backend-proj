package database

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"www.github.com/ic-ETITE-24/icetite-24-backend/config"
)

type RedisRepository struct {
	client *redis.Client
}

var RedisClient *RedisRepository

func NewRepository(redisConfig config.RedisConfig) error {

	client := redis.NewClient(&redis.Options{
		Addr: redisConfig.REDIS_HOST,
		DB:   0,
	})

	// Verify the connection to Redis
	if err := client.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Redis Init Failed: " + err.Error())
		return err
	}

	RedisClient = &RedisRepository{client}
	return nil
}

func (r *RedisRepository) Set(key, value string, time time.Duration) error {
	ctx := context.Background()
	err := r.client.Set(ctx, key, value, time).Err()
	return err
}

func (r *RedisRepository) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisRepository) Delete(key string) error {
	ctx := context.Background()
	err := r.client.Del(ctx, key).Err()
	if err == redis.Nil {
		return nil
	}
	return err
}
