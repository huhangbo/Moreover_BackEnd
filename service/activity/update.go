package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"encoding/json"
)

func UpdateActivity(activity dao.Activity) int {
	if err := conn.MySQL.Model(dao.Activity{}).Where("activity_id = ? AND publisher = ?", activity.ActivityId, activity.Publisher).Updates(activity).Error; err != nil {
		return response.FAIL
	}
	if err := conn.MySQL.Model(&dao.Activity{}).Where("activity_id = ?", activity.ActivityId).First(&activity).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	postJson, _ := json.Marshal(activity)
	if _, err := conn.Redis.Set(key, string(postJson), activityExpiration).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
