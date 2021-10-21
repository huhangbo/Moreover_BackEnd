package follow

import (
	"Moreover/conn"
	"Moreover/dao"
	"github.com/go-redis/redis"
)

func SyncFollowToRedis(follower, category, tmp string) {
	var follows []dao.Follow
	if err := conn.MySQL.Model(&dao.Follow{}).Select(tmp, "created_at").Where(category+" = ?", follower).Find(&follows).Error; err != nil {
		return
	}
	key := category + ":sort:" + follower
	pipe := conn.Redis.Pipeline()
	if category == "publisher" {
		for _, item := range follows {
			pipe.ZAdd(key, redis.Z{Member: item.Parent, Score: float64(item.CreatedAt.Unix())})
		}
	} else {
		for _, item := range follows {
			pipe.ZAdd(key, redis.Z{Member: item.Publisher, Score: float64(item.CreatedAt.Unix())})
		}
	}
	pipe.Expire(key, timeFollowExpiration)
	if _, err := pipe.Exec(); err != nil {
	}
}
