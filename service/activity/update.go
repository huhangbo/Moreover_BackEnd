package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func UpdateActivity(activity dao.Activity, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: activity.ActivityId,
	}
	GetActivityById(&tmpActivity)
	if tmpActivity.Publisher != stuId {
		return response.AuthError
	}
	if err := conn.MySQL.Model(&activity).Updates(activity).Error; err != nil {
		return response.ERROR
	}
	conn.MySQL.Where("activity_id = ?", activity.ActivityId).First(&activity)
	publishActivityToRedis(activity)
	return response.SUCCESS
}
