package cache

import (
	"context"
	"github.com/ohwin/micro/core/config"
	"github.com/ohwin/micro/core/store"

	"github.com/redis/go-redis/v9"
)

func Init() {
	c := config.App.Cache
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Host,
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
	})

	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}

	store.RDB = rdb
}
