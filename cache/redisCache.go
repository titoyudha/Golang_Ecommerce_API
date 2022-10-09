package cache

import (
	"context"
	"encoding/json"
	"go_ecommerce/models"
	"time"

	"github.com/go-redis/redis/v9"
)

type RedisCache struct {
	host    string
	db      int
	expires time.Duration
}

func (cache *RedisCache) NewRedisCache(host string, db int, exp time.Duration) UserCache {
	return &RedisCache{
		host:    host,
		db:      db,
		expires: exp,
	}
}

func (cache *RedisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *RedisCache) Set(key string, value *models.User) {
	client := cache.getClient()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Set(ctx, key, string(json), cache.expires*time.Second)
}

func (cache *RedisCache) Get(key string) *models.User {
	client := cache.getClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	user := models.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		panic(err)
	}
	return &user
}
