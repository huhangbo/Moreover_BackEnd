package follow

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
	"time"
)

func SyncFollowMysqlToRedis(follower, category string) int {
	key := category + ":sort:" + follower
	if category == "follower" {
		category = "fan"
	} else {
		category = "follower"
	}
	var tmpFollows []model.Follow
	sql := `SELECT fan, update_time, follower
			FROM follow
			WHERE ` + category + ` = ?
			AND deleted = 0
			ORDER BY update_time DESC`
	if err := mysql.DB.Select(&tmpFollows, sql, follower); err != nil {
		return response.ERROR
	}
	pipe := redis.DB.Pipeline()
	if category == "follower" {
		for _, item := range tmpFollows {
			tmpPublishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.UpdateTime, time.Local)
			pipe.ZAdd(key, goRedis.Z{Member: item.Fan, Score: float64(tmpPublishTime.Unix())})
		}
	}
	for _, item := range tmpFollows {
		tmpPublishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.UpdateTime, time.Local)
		pipe.ZAdd(key, goRedis.Z{Member: item.Follower, Score: float64(tmpPublishTime.Unix())})
	}
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
