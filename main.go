package main

import (
	"Moreover/conn"
	"Moreover/pkg/mysql"
	"Moreover/router"
	"Moreover/setting"
	"fmt"
)

func main() {
	setting.Init()
	mysql.Init(setting.Config.MySQLConfig)
	conn.InitMysql(setting.Config.MySQLConfig)
	defer mysql.Close()
	conn.InitRedis(setting.Config.RedisConfig)
	defer conn.Close()
	router.InitRouter(fmt.Sprintf(":%d", setting.Config.Port))
}
