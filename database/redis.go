package database

import "github.com/go-redis/redis"

var RDB *redis.Client

func InitRedis(dsn string) *redis.Client {
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opt)
}
