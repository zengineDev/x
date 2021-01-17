package intermittentx

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"

	"github.com/zengineDev/x/configx"
)

type Driver struct {
	Con *redis.Client
}

var (
	instance *Driver
)

var once sync.Once

func GetRedisConnection() *Driver {
	conf := configx.GetConfig()
	once.Do(func() {
		rdb := redis.NewClient(&redis.Options{
			Addr:      fmt.Sprintf("%s:%v", conf.Redis.Host, conf.Redis.Port),
			Password:  conf.Redis.Password, // no password set
			Username:  conf.Redis.User,
			DB:        0,
			TLSConfig: &tls.Config{InsecureSkipVerify: true},
		})

		rdb.Ping(context.Background())
		instance = &Driver{Con: rdb}
	})

	return instance

}
