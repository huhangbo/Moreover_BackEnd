package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func DeleteActivity(activity dao.Activity) int {
	tmpActivity := dao.Activity{ActivityId: activity.ActivityId}
	if code := GetActivityById(&tmpActivity); code != response.SUCCESS {
		return code
	}
	if err := conn.MySQL.Where("activity_id = ? AND publisher = ?", activity.ActivityId, activity.Publisher).Delete(&dao.Activity{}).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	pipe := conn.Redis.Pipeline()
	pipe.ZRem(sortKey, activity.ActivityId)
	pipe.ZRem(sortKey+tmpActivity.Category, activity.ActivityId)
	pipe.Del(key)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
