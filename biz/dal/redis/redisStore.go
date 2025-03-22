package redis

import (
	"context"
	"time"
)

func RedisSet(ctx context.Context, id string, value string) error {
	expire, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}
	return RedisClient.Set(ctx, id, value, expire).Err()
}

func RedisGet(ctx context.Context, id string, clear bool) string {
	val, err := RedisClient.Get(ctx, id).Result()
	if err != nil {
		return ""
	}
	if clear {
		RedisClient.Del(ctx, id)
	}
	return val
}

func RedisVerify(ctx context.Context, id, answer string, clear bool) bool {
	val := RedisGet(ctx, id, clear)
	return val == answer
}
