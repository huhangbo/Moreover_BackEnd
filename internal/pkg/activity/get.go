package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	goRedis "github.com/go-redis/redis"
)

func GetActivityById(activityId string) (int, model.Activity) {
	code, tmpActivity := getActivityByIdFromRedis(activityId)
	if code != response.SUCCESS {
		code, tmpActivity = getActivityByIdFromMysql(activityId)
		if code != response.SUCCESS {
			return code, tmpActivity
		}
		publishActivityToRedis(tmpActivity)
	}
	return code, tmpActivity
}

func GetActivityIdsByPageFromRedis(current, size int, category string) (int, []string) {
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
		fmt.Printf("get activitiesId from redis fail, err:%v\n", errRedis)
		return response.ERROR, nil
	}
	if len(activitiesId) == 0 {
		return response.NotFound, nil
	}
	return response.SUCCESS, activitiesId
}

func GetActivitiesByPageFromMysql(current, size int, category string) (int, []model.Activity) {
	var activities []model.Activity
	if category == "" {
		sql := `SELECT * FROM activity
			WHERE deleted = 0
			ORDER BY update_time
			LIMIT ? ,?`
		err := mysql.DB.Select(activities, sql, (current-1)*size, size)
		if err != nil {
			fmt.Printf("get activities by page fail, err: %v\n", err)
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
		fmt.Printf("get activities by page fail, err: %v\n", err)
		return response.ERROR, nil
	}
	return response.SUCCESS, activities
}

func GetActivityByIds(activityIds []string) (int, []model.Activity) {
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

func GetPublisherById(activityId string) (int, string) {
	code, activity := GetActivityById(activityId)
	if code != response.SUCCESS {
		return code, ""
	}
	return response.SUCCESS, activity.Publisher
}

func GetTotal(category string, size int) (int, int, int) {
	code, total, totalPage := getTotalFromRedis(category, size)
	if code != response.SUCCESS {
		code, total, totalPage = getTotalFromMysql(category, size)
		return code, total, totalPage
	}
	return response.SUCCESS, total, totalPage
}

func getActivityByIdFromRedis(activityId string) (int, model.Activity) {
	var activity model.Activity
	activityString, err := redis.DB.Get("activity:id:" + activityId).Result()
	if err != nil { //err判断
		fmt.Printf("get activity by id from redis fail, err: %v\n", err)
		if err == goRedis.Nil {
			return response.NotFound, activity
		}
		return response.ERROR, activity
	}
	if err := json.Unmarshal([]byte(activityString), &activity); err != nil {
		fmt.Printf("activityString to struct fail, err: %v\n", err)
		return response.ERROR, activity
	}
	return response.SUCCESS, activity
}

func getTotalFromRedis(category string, size int) (int, int, int) {
	key := "activity:sort"
	if category != "" {
		key = "activity:sort:" + category
	}
	total, err := redis.DB.ZCard(key).Result()
	if err != nil {
		fmt.Printf("get activity total fail, err %v\n", err)
		return response.NotFound, 0, 0
	}
	totalPage := int(total)/size + 1
	return response.SUCCESS, int(total), totalPage
}

func getActivityByIdFromMysql(activityId string) (int, model.Activity) {
	var activity model.Activity
	sql := `SELECT * FROM activity
			WHERE activity_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&activity, sql, activityId); err != nil {
		fmt.Printf("get activity by id from mysql fail, err: %v\n", err)
		return response.ERROR, activity
	}
	return response.SUCCESS, activity
}

func getTotalFromMysql(category string, size int) (int, int, int) {
	var total int
	if category == "" {
		sql := `SELECT COUNT(*)
				FROM activity
				WHERE deleted = 0`
		if err := mysql.DB.Get(&total, sql); err != nil {
			fmt.Printf("get activity total from mysql fail, err: %v\n", err)
			return response.ERROR, 0, 0
		}
		totalPage := total/size + 1
		return response.SUCCESS, total, totalPage
	}
	sql := `SELECT COUNT(*)
			FROM activity
			WHERE deleted = 0
			AND category = ?`
	if err := mysql.DB.Get(&total, sql, category); err != nil {
		fmt.Printf("get activity total from mysql fail, err: %v\n", err)
		return response.ERROR, 0, 0
	}
	totalPage := total/size + 1
	return response.SUCCESS, total, totalPage
}
