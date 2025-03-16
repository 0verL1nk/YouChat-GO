package redis

import (
	"context"
	"time"
)

type Store interface {
	// Set sets the digits for the captcha id.
	Set(id string, value string) error

	// Get returns stored digits for the captcha id. Clear indicates
	// whether the captcha must be deleted from the store.
	Get(id string, clear bool) string

	//Verify captcha's answer directly
	Verify(id, answer string, clear bool) bool
}

type RedisStore struct {
	ctx context.Context
}

func (r *RedisStore) Set(id string, value string) error {
	expire, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}
	return RedisClient.Set(r.ctx, id, value, expire).Err()
}

func (r *RedisStore) Get(id string, clear bool) string {
	val, err := RedisClient.Get(r.ctx, id).Result()
	if err != nil {
		return ""
	}
	if clear {
		RedisClient.Del(r.ctx, id)
	}
	return val
}

func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, clear)
	return val == answer
}
