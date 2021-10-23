package util

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"github.com/go-redis/redis"
	"time"
)

func GetTopScore(star int, createdAt time.Time) float64 {
	hour := time.Now().Hour() - createdAt.Hour()
	return float64((star+1)/(hour+2) ^ 2)
}

func TopPost(post dao.PostDetail) int {
	key := "post:sort:"
	score := GetTopScore(post.Star, post.CreatedAt)
	if err := conn.Redis.ZAdd(key, redis.Z{Member: post.PostId, Score: score}); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
