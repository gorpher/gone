package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client redis.Cmdable
}

type RedisOptions struct {
	Address  []string `json:"address" yaml:"address"`
	Username string   `json:"username" yaml:"username"`
	Password string   `json:"password" yaml:"password"`
	DB       int      `json:"db" yaml:"db"`
}

func NewRedisCache(options *RedisOptions) (*RedisCache, error) {
	var client redis.Cmdable
	if len(options.Address) == 1 {
		client = redis.NewClient(&redis.Options{
			Addr:        options.Address[0],
			DB:          options.DB,
			Password:    options.Password,
			ReadTimeout: 5 * time.Second,
		})
	}
	if len(options.Address) > 1 {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    options.Address,
			Username: options.Username,
			Password: options.Password,
		})
	}
	if client != nil {
		if err := client.Ping(context.TODO()).Err(); err != nil {
			return nil, err
		}
	}
	return NewRedisCacheDB(client)
}

func NewRedisCacheDB(client redis.Cmdable) (*RedisCache, error) {
	return &RedisCache{client: client}, nil
}

func (s *RedisCache) Get(key string) ([]byte, error) {
	return s.client.Get(context.Background(), key).Bytes()
}

func (s *RedisCache) Set(key string, value interface{}) error {
	return s.client.Set(context.Background(), key, value, 0).Err()
}

func (s *RedisCache) SetWithTTL(key string, value string, duration time.Duration) error {
	return s.client.SetEx(context.Background(), key, value, duration).Err()
}

func (s *RedisCache) Del(key string) error {
	return s.client.Del(context.Background(), key).Err()
}
