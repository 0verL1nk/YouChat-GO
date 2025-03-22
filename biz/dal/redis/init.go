package redis

import (
	"context"

	"core/conf"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func Init() (err error) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		return err
	}
	return
}
