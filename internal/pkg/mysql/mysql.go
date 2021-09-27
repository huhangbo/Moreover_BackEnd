package mysql

import (
	"Moreover/setting"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Init(config *setting.MySQLConfig)  {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true", config.User, config.Password, config.Host, config.Port, config.DB)
	fmt.Printf(dsn)
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("Connect MySQL falied, err: %v\n", err)
		return
	}
	DB.SetMaxOpenConns(config.MaxIdleConns)
	DB.SetMaxIdleConns(config.MaxIdleConns)
}

func Close()  {
	err := DB.Close()
	if err != nil {
		fmt.Printf("MySQL close failed, err: %v\n", err)
	}
}