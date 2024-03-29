package redis

import (
	"bluebell/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(conf *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Password,
		DB:       conf.DB,
		PoolSize: conf.Pool_Size,
	})
	_, err = client.Ping().Result()
	if err != nil {
		fmt.Printf("connect redis failed, err: %v\n", err)
		return
	}
	return
}

func Close() {
	_ = client.Close()
}
