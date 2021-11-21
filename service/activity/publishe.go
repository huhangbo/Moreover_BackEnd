package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

const (
	activityExpiration = time.Hour * 7 * 24
	sortKey            = "activity:sort:"
)

func PublishActivity(activity dao.Activity) int {
	if err := conn.MySQL.Create(&activity).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	postJson, _ := json.Marshal(activity)
	pipe := conn.Redis.Pipeline()
	pipe.ZAdd(sortKey, redis.Z{Member: activity.ActivityId, Score: float64(activity.CreatedAt.Unix())})
	pipe.ZAdd(sortKey+activity.Category, redis.Z{Member: activity.ActivityId, Score: float64(activity.CreatedAt.Unix())})
	pipe.Set(key, string(postJson), activityExpiration)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
