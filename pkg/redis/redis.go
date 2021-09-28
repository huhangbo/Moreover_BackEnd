package redis

import (
	"Moreover/setting"
	"fmt"
	"github.com/go-redis/redis"
)

var DB *redis.Client

func Init(config *setting.RedisConfig) {
	DB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB, // use default DB
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})
	if _, err := DB.Ping().Result(); err != nil {
		fmt.Printf("Connect redis failed, err: %v\n", err)
		panic(err)
	}

}

func Close() {
	err := DB.Close()
	if err != nil {
		fmt.Printf("MySQL close failed, err: %v\n", err)
	}
}
