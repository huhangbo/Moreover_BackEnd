package util

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"github.com/go-redis/redis"
	"time"
)

func GetTopScore(star, comments int, createdAt time.Time) float64 {
	minutes := time.Now().Minute() - createdAt.Minute()
	return (float64(star)*0.3 + float64(comments)*0.7) * 1000 / (float64(minutes)/60 + 1.5)
}

func TopPost(post dao.PostDetail) int {
	var (
		key = "post:sort:top"
		//incr int
	)

	score := GetTopScore(post.Star, post.Comments, post.CreatedAt)
	if err := conn.Redis.ZAdd(key, redis.Z{Member: post.PostId, Score: score}); err != nil {
		return response.FAIL
	}
	//conn.Redis.ZIncrBy(key, )
	return response.SUCCESS
}
