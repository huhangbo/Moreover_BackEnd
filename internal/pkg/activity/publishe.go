package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	goRedis "github.com/go-redis/redis"
	"time"
)

const activityExpiration = time.Hour * 24 * 7

func PublishActivity(activity model.Activity) int {
	code := publishActivityToRedis(activity)
	if code != response.SUCCESS {
		return code
	}
	code = publishActivityToMysql(activity)
	if code != response.SUCCESS {
		deleteActivityFromRedis(activity.ActivityId, activity.Category)
	}
	return code
}

func publishActivityToRedis(activity model.Activity) int {
	jsonActivity, err := json.Marshal(activity)
	publishTime, _ := time.ParseInLocation("2006/01/02 15:05:06", activity.CreateTime, time.Local)
	if err != nil {
		return response.ERROR
	}
	activityKey := "activity:id:" + activity.ActivityId
	sortCategoryKey := "activity:sort:" + activity.Category
	sortKey := "activity:sort"
	sortActivity := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: activity.ActivityId,
	}
	pipe := redis.DB.Pipeline()
	pipe.ZAdd(sortCategoryKey, sortActivity)
	pipe.ZAdd(sortKey, sortActivity)
	pipe.Set(activityKey, string(jsonActivity), activityExpiration)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func publishActivityToMysql(activity model.Activity) int {
	sql := `INSERT INTO activity (create_time, update_time, activity_id, publisher, category, title, outline, start_time, end_time, contact, location, detail) 
			VALUES (:create_time, :update_time, :activity_id, :publisher, :category, :title, :outline, :start_time, :end_time, :contact, :location, :detail)`
	if _, err := mysql.DB.NamedExec(sql, activity); err != nil {
		code, _ := getActivityByIdFromMysql(activity.ActivityId)
		return code
	}
	return response.SUCCESS
}
