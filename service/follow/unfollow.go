package follow

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
)

func Unfollow(follower, fan string) int {
	code := UnfollowRedis(follower, fan)
	if code != response.SUCCESS {
		return code
	}
	code = UnfollowMysql(follower, fan, 1)
	if code != response.SUCCESS {
		PublishFollowToRedis(model.Follow{
			Follower: follower,
			Fan:      fan,
		})
	}
	return code
}

func UnfollowRedis(follower, fan string) int {
	keyFan := "follow:sort:" + follower
	keyFollow := "fan:sort:" + fan
	pipe := redis.DB.Pipeline()
	pipe.ZRem(keyFan, follower)
	pipe.ZRem(keyFollow, fan)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func UnfollowMysql(follower, fan string, state int) int {
	sql := `UPDATE follow 
			SET deleted = ?
			WHERE follower = ?
			AND fan = ?`
	if _, err := mysql.DB.Exec(sql, state, follower, fan); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
