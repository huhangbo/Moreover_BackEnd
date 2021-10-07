package activity

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
)

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
				_, tmp := user.GetUserInfo(activities[i].Publisher)
				activities[i].PublisherInfo = tmp.UserBasicInfo
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
		err := mysql.DB.Select(&activities, sql, (current-1)*size, size)
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
	err := mysql.DB.Select(&activities, sql, category, (current-1)*size, size)
	if err != nil {
		return response.ERROR, nil
	}
	return response.SUCCESS, activities
}
