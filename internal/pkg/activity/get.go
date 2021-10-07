package activity

import (
	"Moreover/internal/pkg/user"
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
	_, tmpActivity.Star = util.GetTotalById(tmpActivity.ActivityId, "liked")
	return code, tmpActivity
}

func GetActivitiesByPade(current, size int, category string) (int, []model.Activity, model.Page) {
	var activities []model.Activity
	code, total := GetTotal(category)
	var tmpPage = model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: (total / size) + 1,
	}
	if code != response.SUCCESS {
		return code, activities, tmpPage
	}
	if (current-1)*size > total {
		return response.ParamError, activities, tmpPage
	}
	codeIdsRedis, activityIds := getActivityIdsByPageFromRedis(current, size, category)
	if codeIdsRedis != response.SUCCESS || (len(activityIds) == 0 && code == response.SUCCESS) {
		code, activities = getActivitiesByPageFromMysql(current, size, category)
		if code == response.SUCCESS {
			for i := 0; i < len(activities); i++ {
				PublishActivity(activities[i])
				_, activities[i].PublisherInfo = user.GetUserInfo(activities[i].Publisher)
				activities[i].PublisherInfo.Description = ""
			}
		}
		return code, activities, tmpPage
	}
	code, activities = getActivityByIds(activityIds)
	if code != response.SUCCESS {
		return code, activities, tmpPage
	}
	return response.SUCCESS, activities, tmpPage
}

func GetPublisherById(activityId string) (int, string) {
	code, activity := GetActivityById(activityId)
	if code != response.SUCCESS {
		return code, ""
	}
	return response.SUCCESS, activity.Publisher
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
		_, tmpRedisActivity.PublisherInfo = user.GetUserInfo(tmpRedisActivity.Publisher)
		tmpRedisActivity.PublisherInfo.Description = ""
		activities = append(activities, tmpRedisActivity)
	}
	return response.SUCCESS, activities
}

func getActivityIdsByPageFromRedis(current, size int, category string) (int, []string) {
	sortKey := "activity:sort"
	if category != "" {
		sortKey = "activity:sort:" + category
	}
	rangeOpt := goRedis.ZRangeBy{
		Min:    "-",
		Max:    "+",
		Offset: int64((current - 1) * size),
		Count:  int64(size),
	}
	activitiesId, errRedis := redis.DB.ZRangeByLex(sortKey, rangeOpt).Result()
	if errRedis != nil {
		return response.ERROR, nil
	}
	if len(activitiesId) == 0 {
		return response.NotFound, nil
	}
	return response.SUCCESS, activitiesId
}

func getActivitiesByPageFromMysql(current, size int, category string) (int, []model.Activity) {
	var activities []model.Activity
	if category == "" {
		sql := `SELECT * FROM activity
			WHERE deleted = 0
			ORDER BY update_time
			LIMIT ? ,?`
		err := mysql.DB.Select(activities, sql, (current-1)*size, size)
		if err != nil {
			return response.ERROR, nil
		}
		return response.SUCCESS, activities
	}
	sql := `SELECT * FROM activity
			WHERE category = ?
			AND deleted = 0
			ORDER BY update_time
			LIMIT ? ,?`
	err := mysql.DB.Select(activities, sql, category, (current-1)*size, size)
	if err != nil {
		return response.ERROR, nil
	}
	return response.SUCCESS, activities
}

func getActivityByIdFromRedis(activityId string) (int, model.Activity) {
	var activity model.Activity
	activityString, err := redis.DB.Get("activity:id:" + activityId).Result()
	if err != nil { //err判断
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
	total, err := redis.DB.ZCard(key).Result()
	if err != nil {
		return response.NotFound, 0
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
	if err := mysql.DB.Select(&total, sql, category); err != nil {
		return response.ERROR, 0
	}
	return response.SUCCESS, total
}
