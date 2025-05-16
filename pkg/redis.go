package pkg

import (
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis(redisAddr string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
}

func GetClient() *redis.Client {
	return Rdb
}
