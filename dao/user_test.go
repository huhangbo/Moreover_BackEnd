package dao

import (
	"Moreover/connent"
	"Moreover/setting"
	"testing"
)

var tmp = UserInfo{
	StudentId:   "17195477",
	Nickname:    "秘密",
	Avatar:      "sss",
	Sex:         "男",
	Description: "顶顶顶",
}

func TestUser_Add(t *testing.T) {
	setting.Init()
	connent.InitMysql(setting.Config.MySQLConfig)
	connent.MySQL.AutoMigrate(tmp)
	connent.MySQL.Create(&tmp)
}

func TestUser_Get(t *testing.T) {
	setting.Init()
	connent.InitMysql(setting.Config.MySQLConfig)
}
