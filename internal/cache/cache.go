package cache

import (
	entity2 "AvitoWinter/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClientAdapter struct {
	client *redis.Client
}

func NewRedisClientAdapter(ctx context.Context) (*RedisClientAdapter, error) {
	var redisConf RedisConfig
	err := redisConf.UpdateEnvAddress()
	if err != nil {
		return nil, fmt.Errorf("-> redisConf.UpdateEnvAddress%v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConf.Addr,
		Password: redisConf.Passwd,
		DB:       redisConf.DB,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("-> rdb.Ping%v", err)
	}

	client := &RedisClientAdapter{client: rdb}

	return client, nil
}

func (r *RedisClientAdapter) Set(ctx context.Context, key string, user *entity2.User, expiration time.Duration) error {
	value, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("-> json.Marshal:%v", err)
	}

	err = r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("-> r.client.Set: %v", err)
	}

	return nil
}

func (r *RedisClientAdapter) Get(ctx context.Context, key string) (*entity2.User, error) {
	value, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("->  r.client.Get.Result: %v", err)
	}

	var user entity2.User
	err = json.Unmarshal(value, &user)
	if err != nil {
		return nil, fmt.Errorf("CreateUser-> json.NewDecoder: неверный формат для пользователя: %v", err)
	}

	return &user, nil
}

//type Client interface {
//	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
//	Get(ctx context.Context, key string) error
//}
//
//type Cache struct {
//	cacheClient Client
//}
//
//func NewCacheClient(client Client) *Cache {
//	return &Cache{cacheClient: client}
//}
