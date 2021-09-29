package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	goRedis "github.com/go-redis/redis"
)

func GetActivityByPage(current, size int) []model.Activity {
	var activities []model.Activity
	sql := `SELECT * FROM activity
			WHERE deleted = 0 
			LIMIT ? ,?`
	err := mysql.DB.Select(&activities, sql, (current-1)*size, size)
	if err != nil {
		fmt.Printf("get activities by page fail, err: %v\n", err)
		return nil
	}
	return activities
}

func GetActivityById(activityId string, activity *model.Activity) int {
	activityString, err := redis.DB.Get("activity:id:" + activityId).Result()
	if err != goRedis.Nil && err != nil {
		fmt.Printf("get activity by id from redis fail, err: %v\n", err)
		return response.ERROR
	}
	if err == goRedis.Nil {
		if err := json.Unmarshal([]byte(activityString), activity); err != nil {
			fmt.Printf("activityString to struct fail, err: %v\n", err)
			return response.ERROR
		}
	}
	sql := `SELECT * FROM activity
			WHERE activity_id = ?`
	if err := mysql.DB.Get(activity, sql, activityId); err != nil {
		fmt.Printf("get activity by id from mysql fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}
