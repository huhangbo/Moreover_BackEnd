package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"encoding/json"
)

func GetActivityById(activity *dao.Activity) int {
	code := getActivityByIdFromRedis(activity)
	if code != response.SUCCESS {
		if err := conn.MySQL.Where("activity_id = ?", (*activity).ActivityId).First(activity).Error; err != nil {
			return response.FAIL
		}
		publishActivityToRedis(*activity)
		return response.SUCCESS
	}
	return code
}

func GetActivityDetailById(detail *dao.ActivityDetail, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: detail.ActivityId,
	}
	code := GetActivityById(&tmpActivity)
	detail.Activity = tmpActivity
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	_, detail.Star, detail.IsStar = liked.GetTotalAndLiked(detail.ActivityId, stuId)
	return code
}

func GetTotal(category string) (int, int) {
	key := "activity:sort:" + category
	total, err := conn.Redis.ZCard(key).Result()
	if err != nil || total == 0 {
		go SyncActivitySortRedis()
		if category == "" {
			if err := conn.MySQL.Model(&dao.Activity{}).Count(&total).Error; err != nil {
				return response.ERROR, int(total)
			}
		} else {
			if err := conn.MySQL.Model(&dao.Activity{}).Where("category = ?", category).Count(&total); err != nil {
				return response.ERROR, int(total)
			}
		}
		return response.SUCCESS, int(total)
	}
	return response.SUCCESS, int(total)
}

func getActivityByIds(activityIds []string) (int, []dao.Activity) {
	var activities []dao.Activity
	for i := 0; i < len(activityIds); i++ {
		tmpActivity := dao.Activity{
			ActivityId: activityIds[i],
		}
		GetActivityById(&tmpActivity)
		activities = append(activities, tmpActivity)
	}
	return response.SUCCESS, activities
}

func getActivityByIdFromRedis(activity *dao.Activity) int {
	activityString, err := conn.Redis.Get("activity:id:" + (*activity).ActivityId).Result()
	if activityString == "" || err != nil {
		return response.FAIL
	}
	if err := json.Unmarshal([]byte(activityString), activity); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
