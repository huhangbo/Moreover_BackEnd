package activity

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
)

func GetActivitiesByPade(current, size int, category string) (int, []model.ActivityPageShow, model.Page) {
	var activities []model.ActivityPageShow
	var tmpActivities []model.Activity
	var activityIds []string
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
	code, activityIds = getActivityIdsByPageFromRedis(current, size, category)
	code, tmpActivities = getActivityByIds(activityIds)
	if code != response.SUCCESS || (len(activityIds) == 0 && code == response.SUCCESS) {
		code, tmpActivities = getActivitiesByPageFromMysql(current, size, category)
		if code == response.SUCCESS {
			go SyncActivitySortRedisMysql()
		}
	}
	for i := 0; i < len(tmpActivities); i++ {
		var tmpPageShow model.ActivityPageShow
		_, tmpUser := user.GetUserInfo(tmpActivities[i].Publisher)
		tmpPageShow.ActivityBasic = tmpActivities[i].ActivityBasic
		tmpPageShow.PublisherInfo = tmpUser.UserBasicInfo
		activities = append(activities, tmpPageShow)
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
