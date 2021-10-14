package main

import (
	"Moreover/connent"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/router"
	"Moreover/setting"
	"fmt"
)

func main() {
	setting.Init()
	mysql.Init(setting.Config.MySQLConfig)
	connent.InitMysql(setting.Config.MySQLConfig)
	defer mysql.Close()
	redis.Init(setting.Config.RedisConfig)
	defer redis.Close()
	router.InitRouter(fmt.Sprintf(":%d", setting.Config.Port))
}
