package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"online-store/pkg/config"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(conf *config.Config) *Cache {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Pass,
		DB:       conf.Redis.DB,
	})

	return &Cache{rdb: rdb}
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		slog.Error(
			"Failed to set cache",
			"key", key,
			"value", value,
			"expiration", expiration,
		)
		return err
	}

	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}

		slog.Error(
			"Failed to get cache",
			"key", key,
			"err", err,
		)
		return "", err
	}

	return val, nil
}

func (c *Cache) Del(ctx context.Context, key string) error {
	keys := []string{}

	if strings.Contains(key, "*") {
		keyList := c.rdb.Keys(ctx, key)
		keys = append(keys, keyList.Val()...)
	} else {
		keys = []string{key}
	}

	for _, k := range keys {
		err := c.rdb.Del(ctx, k).Err()
		if err != nil {
			slog.Error(
				"Failed to delete cache",
				"key", key,
			)
			return err
		}
	}

	return nil
}
