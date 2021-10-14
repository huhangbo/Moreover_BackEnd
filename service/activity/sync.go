package activity

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	goRedis "github.com/go-redis/redis"
	"time"
)

func SyncActivitySortRedisMysql() {
	var activities []model.Activity
	sql := `SELECT * FROM activity
			WHERE deleted = 0
			ORDER BY update_time DESC`
	if err := mysql.DB.Select(&activities, sql); err != nil {
		return
	}
	pipe := redis.DB.Pipeline()
	for _, item := range activities {
		key := "activity:sort:" + item.Category
		tmpPublishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.UpdateTime, time.Local)
		pipe.ZAdd(key, goRedis.Z{Member: item.ActivityId, Score: float64(tmpPublishTime.Unix())})
		pipe.ZAdd("activity:sort:", goRedis.Z{Member: item.ActivityId, Score: float64(tmpPublishTime.Unix())})
	}
	if _, err := pipe.Exec(); err != nil {
	}
	return
}
