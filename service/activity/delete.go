package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func DeleteActivity(activity dao.Activity, stuId string) int {
	code := GetActivityById(&activity)
	if code != response.SUCCESS {
		return response.ERROR
	}
	if activity.Publisher != stuId {
		return response.AuthError
	}
	conn.MySQL.Delete(&activity)
	code = deleteActivityFromRedis(activity)
	return code
}

func deleteActivityFromRedis(activity dao.Activity) int {
	keyActivity := "activity:id:" + activity.ActivityId
	keySort := "activity:sort:"
	keyCategorySort := "activity:sort:" + activity.Category
	pipe := conn.Redis.Pipeline()
	pipe.Del(keyActivity)
	pipe.ZRem(keySort, activity.ActivityId)
	pipe.ZRem(keyCategorySort, activity.ActivityId)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
