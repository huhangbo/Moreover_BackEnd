package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
)

func DeleteActivityById(studentId, activityId string) int {
	activityString, err := redis.DB.Get("activity:id:" + activityId).Result()
	if err != nil {
		fmt.Printf("get activity from redis fail, err: %v\n", err)
		return response.ERROR
	}
	var activityStruct model.Activity
	if err := json.Unmarshal([]byte(activityString), &activityStruct); err != nil {
		fmt.Printf("activity json to struct fail, err: %v\n", err)
		return response.ERROR
	}
	if activityStruct.Publisher != studentId {
		return response.AuthError
	}
	sql := `UPDATE activity
			SET deleted = 1
			WHERE activity_id = ?`
	if _, err := mysql.DB.Exec(sql, activityId); err != nil {
		fmt.Printf("delete activity from mysql fail, err: %v\n", err)
		return response.ERROR
	}
	activityStruct.Deleted = 1
	activityJson, _ := json.Marshal(activityStruct)
	redis.DB.Set("activity:id:"+activityId, activityJson, activityExpiration)
	return response.SUCCESS
}
