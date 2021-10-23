package main

import (
	"Moreover/conn"
	"Moreover/router"
	"Moreover/setting"
	"fmt"
)

func main() {
	setting.Init()

	conn.InitMysql(setting.Config.MySQLConfig)

	conn.InitRedis(setting.Config.RedisConfig)

	defer conn.RedisClose()

	conn.InitMongo(setting.Config.MongoConfig)

	router.InitRouter(fmt.Sprintf(":%d", setting.Config.Port))
}
