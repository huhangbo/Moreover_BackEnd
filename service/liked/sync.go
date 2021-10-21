package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"github.com/go-redis/redis"
)

func SyncLikeToRedis(likes []dao.Liked) {
	key := "like:sort:" + likes[0].ParentId
	pipe := conn.Redis.Pipeline()
	for _, item := range likes {
		pipe.ZAdd(key, redis.Z{Member: item.Publisher, Score: float64(item.CreatedAt.Unix())})
	}
	pipe.Expire(key, timeLikedExpiration)
	if _, err := pipe.Exec(); err != nil {
		return
	}
	return
}
