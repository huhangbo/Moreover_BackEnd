package redis

import (
	"Moreover/setting"
	"fmt"
	"github.com/go-redis/redis"
)

var client *redis.Client

func Init (config *setting.RedisConfig) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password, // no password set
		DB:           config.DB,       // use default DB
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
	})
	if _, err := client.Ping().Result(); err != nil {
		fmt.Printf("Connect redis failed, err: %v\n", err)
	}

}

func Close()  {
	err := client.Close()
	if err != nil {
		fmt.Printf("MySQL close failed, err: %v\n", err)
	}
}