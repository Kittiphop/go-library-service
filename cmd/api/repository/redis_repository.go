package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type RedisRepository struct {
	redis *redis.Client
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func NewRedisRepository(config RedisConfig) (*RedisRepository, error){
	redis := redis.NewClient(&redis.Options{
		Addr:     config.Host+":"+config.Port,
		Password: config.Password,
		DB:       0,
	})

	_, err := redis.Ping(context.Background()).Result()
	if err != nil {
		return nil, errors.Wrap(err, "[NewRedisRepository]: unable to connect redis")	
	}

	return &RedisRepository{redis: redis}, nil
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.redis.Get(context.Background(), key).Result()
}

func (r *RedisRepository) Set(key string, value interface{}, expiration uint) error {
	return r.redis.Set(context.Background(), key, value, time.Duration(expiration)*time.Second).Err()
}

func (r *RedisRepository) Delete(key string) error {
	return r.redis.Del(context.Background(), key).Err()
}

