package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	"time"
)

const activityExpiration = time.Hour * 24 * 7

func PublishActivity(activity model.Activity) int {
	sql := `INSERT INTO activity (create_time, update_time, activity_id, publisher, category, title, outline, start_time, end_time, contact, location, detail) 
			VALUES (:create_time, :update_time, :activity_id, :publisher, :category, :title, :outline, :start_time, :end_time, :contact, :location, :detail)`
	if _, err := mysql.DB.NamedExec(sql, activity); err != nil {
		fmt.Printf("insert activity to mysql fali, err: %v\n", err)
		return response.ParamError
	}
	jsonActivity, err := json.Marshal(activity)
	if err != nil {
		fmt.Printf("activityStruct to json fail, err: %v\n", err)
		return response.ERROR
	}
	pipe := redis.DB.Pipeline()
	pipe.SetNX("activity:id:"+activity.ActivityId, string(jsonActivity), activityExpiration)
	pipe.Incr("activity:count")
	if _, err := redis.DB.Set("activity:id:"+activity.ActivityId, string(jsonActivity), activityExpiration).Result(); err != nil {
		fmt.Printf("activityStruct to redis fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}
