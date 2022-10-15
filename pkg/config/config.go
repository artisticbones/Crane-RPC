package config

import (
	"github.com/artisticbones/Crane-RPC/database"
	"github.com/go-redis/redis"
)

type Runtime interface {
	SetRDB(url string)
	GetRDB() *redis.Client
}

type Application struct {
	rdb *redis.Client
}

func (a *Application) SetRDB(url string) {
	a.rdb = database.InitRedis(url)
}

func (a *Application) GetRDB() *redis.Client {
	return a.rdb
}
