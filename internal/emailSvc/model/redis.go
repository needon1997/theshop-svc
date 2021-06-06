package model

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	config2 "github.com/needon1997/theshop-svc/internal/common/config"
	"time"
)

func GetCodeByEmail(email string) (string, error) {
	c := newClient()
	return c.Get(context.Background(), email).Result()
}

func SaveEmailCode(email string, code string) error {
	c := newClient()
	return c.Set(context.Background(), email, code, time.Duration(config2.ServerConfig.RedisConfig.Expiration)).Err()
}

func DeleteEmailCodeByEmail(email string) error {
	c := newClient()
	return c.Del(context.Background(), email).Err()
}

func newClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config2.ServerConfig.RedisConfig.Host, config2.ServerConfig.RedisConfig.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return rdb
}
