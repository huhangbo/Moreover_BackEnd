package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
	"time"
)

const activityExpiration = time.Hour * 24 * 7

func PublishActivity(activity dao.Activity) int {
	if err := conn.MySQL.Create(&activity).Error; err != nil {
		return response.FAIL
	}
	sortCategoryKey := "activity:sort:" + activity.Category
	sortKey := "activity:sort:"
	util.PublishSortRedis(activity.ActivityId, float64(time.Now().Unix()), sortCategoryKey, sortKey)
	return response.SUCCESS
}

func publishActivityToRedis(activity dao.Activity) int {
	activityKey := "activity:id:" + activity.ActivityId
	jsonActivity, _ := json.Marshal(activity)
	conn.Redis.Set(activityKey, string(jsonActivity), activityExpiration)
	return response.SUCCESS
}
