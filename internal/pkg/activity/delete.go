package activity

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"fmt"
)

func DeleteActivityById(activityId, stuId string) int {
	code, tmpActivity := GetActivityById(activityId)
	if code != response.SUCCESS {
		return response.ERROR
	}
	if tmpActivity.Publisher != stuId {
		return response.AuthError
	}
	code = deleteActivityFromRedis(activityId, tmpActivity.Category)
	code = deleteActivityFromMysql(activityId)
	return code
}

func deleteActivityFromRedis(activityId, category string) int {
	keyActivity := "activity:id:" + activityId
	keySort := "activity:sort"
	keyCategorySort := "activity:sort:" + category
	pipe := redis.DB.Pipeline()
	pipe.Del(keyActivity)
	pipe.ZRem(keySort, activityId)
	pipe.ZRem(keyCategorySort, activityId)
	if _, err := pipe.Exec(); err != nil {
		fmt.Printf("delete activity from redis fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}

func deleteActivityFromMysql(activityId string) int {
	sql := `UPDATE activity
			SET deleted = 1
			WHERE activity_id = ?`
	if _, err := mysql.DB.Exec(sql, activityId); err != nil {
		fmt.Printf("delete activity from mysql fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}
