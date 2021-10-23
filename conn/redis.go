package conn

import (
	"Moreover/setting"
	"fmt"
	"github.com/go-redis/redis"
)

var Redis *redis.Client

func InitRedis(config *setting.RedisConfig) {
	Redis = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})
	if _, err := Redis.Ping().Result(); err != nil {
		panic(err)
	}
}

func RedisClose() {
	if err := Redis.Close(); err != nil {
	}
}
