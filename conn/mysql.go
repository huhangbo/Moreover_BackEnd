package conn

import (
	"Moreover/setting"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var MySQL *gorm.DB

func InitMysql(config *setting.MySQLConfig) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local", config.User, config.Password, config.Host, config.Port, config.DB)

	MySQL, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	sqlDB, _ := MySQL.DB()

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)

	sqlDB.SetConnMaxLifetime(time.Hour)

}
