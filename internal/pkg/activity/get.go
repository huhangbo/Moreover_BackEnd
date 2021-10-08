package activity

import (
	"Moreover/internal/util"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	goRedis "github.com/go-redis/redis"
)

func GetActivityById(activityId string) (int, model.Activity) {
	code, tmpActivity := getActivityByIdFromRedis(activityId)
	if code != response.SUCCESS {
		code, tmpActivity = getActivityByIdFromMysql(activityId)
		if code == response.SUCCESS {
			publishActivityToRedis(tmpActivity)
		}
	}
	_, tmpActivity.Star = util.GetTotalById(tmpActivity.ActivityId, "liked", "parent_id")
	return code, tmpActivity
}

func GetPublisherById(activityId string) (int, string) {
	code, activity := GetActivityById(activityId)
	return code, activity.Publisher
}

func GetTotal(category string) (int, int) {
	code, total := getTotalFromRedis(category)
	if code != response.SUCCESS {
		code, total = getTotalFromMysql(category)
	}
	return code, total
}

func getActivityByIds(activityIds []string) (int, []model.Activity) {
	var activities []model.Activity
	for i := 0; i < len(activityIds); i++ {
		tmpRedisCode, tmpRedisActivity := GetActivityById(activityIds[i])
		if tmpRedisCode != response.SUCCESS {
			return tmpRedisCode, activities
		}
		activities = append(activities, tmpRedisActivity)
	}
	return response.SUCCESS, activities
}

func getActivityByIdFromRedis(activityId string) (int, model.Activity) {
	var activity model.Activity
	activityString, err := redis.DB.Get("activity:id:" + activityId).Result()
	if err != nil {
		if err == goRedis.Nil {
			return response.NotFound, activity
		}
		return response.ERROR, activity
	}
	if err := json.Unmarshal([]byte(activityString), &activity); err != nil {
		return response.ERROR, activity
	}
	return response.SUCCESS, activity
}

func getTotalFromRedis(category string) (int, int) {
	key := "activity:sort"
	if category != "" {
		key = "activity:sort:" + category
	}
	total, _ := redis.DB.ZCard(key).Result()
	if total == 0 {
		return response.NotFound, int(total)
	}
	return response.SUCCESS, int(total)
}

func getActivityByIdFromMysql(activityId string) (int, model.Activity) {
	var activity model.Activity
	sql := `SELECT * FROM activity
			WHERE activity_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&activity, sql, activityId); err != nil {
		return response.ERROR, activity
	}
	return response.SUCCESS, activity
}

func getTotalFromMysql(category string) (int, int) {
	var total int
	if category == "" {
		sql := `SELECT COUNT(*)
				FROM activity
				WHERE deleted = 0`
		if err := mysql.DB.Get(&total, sql); err != nil {
			return response.ERROR, 0
		}
		return response.SUCCESS, total
	}
	sql := `SELECT COUNT(*)
			FROM activity
			WHERE deleted = 0
			AND category = ?`
	if err := mysql.DB.Get(&total, sql, category); err != nil {
		return response.ERROR, total
	}
	return response.SUCCESS, total
}
