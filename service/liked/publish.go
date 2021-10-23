package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"github.com/go-redis/redis"
	"time"
)

const timeLikedExpiration = time.Hour * 24 * 7

func PublishLike(liked dao.Liked) int {
	if err := conn.MySQL.Create(&liked).Error; err != nil {
		return response.FAIL
	}
	conn.MySQL.First(&liked)
	if !util.PublishSortRedis(liked.Publisher, float64(liked.CreatedAt.Unix()), "liked:sort:"+liked.ParentId) {
		return response.FAIL
	}
	return response.SUCCESS
}

func PublishTopPost(liked dao.Liked) int {
	totalLike, _ := conn.Redis.ZCard("like:sort:" + liked.ParentId).Result()
	score := util.GetTopScore(int(totalLike), liked.CreatedAt)
	if _, err := conn.Redis.ZAdd("post:top:", redis.Z{Member: liked.ParentId, Score: score}).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
