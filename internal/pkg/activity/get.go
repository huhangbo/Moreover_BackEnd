package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"fmt"
)

func GetActivityByPage(current, size int) []model.Activity {
	var activities []model.Activity
	sql := `SELECT * FROM activity LIMIT ?, ? ORDER BY update_time DESC`
	if err := mysql.DB.Get(&activities, sql, (current-1)*size, size); err != nil {
		fmt.Printf("get activities by page fail, err: %v\n", err)
		panic(err)
		return nil
	}
	return activities
}
