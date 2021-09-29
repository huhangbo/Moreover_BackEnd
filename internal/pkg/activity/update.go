package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
	"fmt"
)

func UpdateActivityById(stuId string, tmp model.Activity) int {
	sqlGetId := `SELECT publisher 
				FROM activity
				WHERE activity_id = ?`
	var tmpStuId string
	if err := mysql.DB.Get(&tmpStuId, sqlGetId, tmp.ActivityId); err != nil {
		fmt.Printf("get publisher from mysql activity fail, err: %v\n", err)
		return response.ERROR
	}
	if tmpStuId != stuId {
		return response.AuthError
	}
	sqlUpdate := `UPDATE activity
				  SET update_time = :update_time, title = :title, category = :category, outline = :outline, start_time =:start_time, end_time = :end_time, contact = :contact, location = :location, detail = :detail
				  WHERE activity_id = :activity_id`
	if _, err := mysql.DB.NamedExec(sqlUpdate, tmp); err != nil {
		fmt.Printf("update activity to mysql fali, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}
