package main

import (
	"Moreover/internal/pkg/mysql"
	"Moreover/internal/pkg/redis"
	"Moreover/router"
	"Moreover/setting"
	"fmt"
)

func main()  {
	setting.Init()
	mysql.Init(setting.Config.MySQLConfig)
	defer mysql.Close()
	redis.Init(setting.Config.RedisConfig)
	defer redis.Close()
	router.InitRouter(fmt.Sprintf(":%d", setting.Config.Port))
}