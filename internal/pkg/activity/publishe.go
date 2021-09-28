package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"encoding/json"
	"fmt"
	"time"
)

const activityExpiration = time.Hour * 24 * 7

func PublishActivity(activity model.Activity) bool {
	sql := `INSERT INTO activity (activity_id, publisher, category, title, outline, start_time, end_time, contact, location, detail) 
			VALUES (:activity_id, :publisher, :category, :title, :outline, :start_time, :end_time, :contact, :location, :detail)`
	if _, err := mysql.DB.NamedExec(sql, activity); err != nil {
		fmt.Printf("insert activity to mysql fali, err: %v\n", err)
		panic(err)
		return false
	}
	jsonActivity, err := json.Marshal(activity)
	if err != nil {
		fmt.Printf("activityStruct to json fail, err: %v\n", err)
		return false
	}
	if _, err := redis.DB.Set("activity:"+activity.ActivityId, string(jsonActivity), activityExpiration).Result(); err != nil {
		fmt.Printf("activityStruct to redis fail, err: %v\n", err)
		return false
	}
	return true
}
