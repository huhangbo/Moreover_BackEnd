package follow

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
	"time"
)

func PublishFollow(follow model.Follow) int {
	code := PublishFollowToRedis(follow)
	if code != response.SUCCESS {
		return code
	}
	code = PublishFollowToMysql(follow)
	if code != response.SUCCESS {
		UnfollowRedis(follow.Follower, follow.Fan)
	}
	return code
}

func PublishFollowToRedis(follow model.Follow) int {
	publishTime, _ := time.ParseInLocation("2006/01/02 15:05:06", follow.UpdateTime, time.Local)
	keyFan := "fan:sort:" + follow.Follower
	keyFollow := "follower:sort:" + follow.Fan
	sortFan := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: follow.Fan,
	}
	sortFollow := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: follow.Follower,
	}
	pipe := redis.DB.Pipeline()
	pipe.ZAdd(keyFollow, sortFollow)
	pipe.ZAdd(keyFan, sortFan)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func PublishFollowToMysql(follow model.Follow) int {
	sql := `INSERT follow (create_time, update_time, follower, fan)
			VALUES (:create_time, :update_time, :follower, :fan)`
	_, err := mysql.DB.NamedExec(sql, follow)
	if err != nil {
		code := UnfollowMysql(follow.Follower, follow.Fan, 0)
		return code
	}
	return response.SUCCESS
}
