package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	goRedis "github.com/go-redis/redis"
	"time"
)

func UpdateActivityById(tmp, old model.Activity) int {
	code, tmpPublisher := GetPublisherById(tmp.ActivityId)
	if code != response.SUCCESS {
		return code
	}
	if tmpPublisher != tmp.Publisher {
		return response.AuthError
	}
	code = updateActivityToMysql(tmp)
	if code != response.SUCCESS {
		return code
	}
	code = updateActivityToRedis(tmp, old)
	return code
}

func updateActivityToMysql(activity model.Activity) int {
	sqlUpdate := `UPDATE activity
				  SET update_time = :update_time, title = :title, category = :category, outline = :outline, start_time =:start_time, end_time = :end_time, contact = :contact, location = :location, detail = :detail
				  WHERE activity_id = :activity_id`
	if _, err := mysql.DB.NamedExec(sqlUpdate, activity); err != nil {
		fmt.Printf("update activity to mysql fali, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}

func updateActivityToRedis(activity, old model.Activity) int {
	activityJson, err := json.Marshal(activity)
	publishTime, _ := time.ParseInLocation("2006/01/02 15:05:06", activity.CreateTime, time.Local)
	if err != nil {
		fmt.Printf("activity struct to json fail, err:%v\n", err)
		return response.ERROR
	}
	key := "activity:id:" + activity.ActivityId
	sortCategoryKey := "activity:sort:" + activity.Category
	oldSortKey := "activity:sort:" + old.Category
	sortActivity := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: activity.ActivityId,
	}
	pipe := redis.DB.Pipeline()
	pipe.Set(key, string(activityJson), activityExpiration)
	pipe.ZRem(oldSortKey, activity.ActivityId)
	pipe.ZAdd(sortCategoryKey, sortActivity)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
