package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"encoding/json"
	"github.com/go-redis/redis"
)

func GetPost(post *dao.Post) int {
	key := "post:id:" + post.PostId
	stringDetail, err := conn.Redis.Get(key).Result()
	if err := json.Unmarshal([]byte(stringDetail), post); err != nil {
		return response.FAIL
	}
	if err != nil {
		if err := conn.MySQL.First(post).Error; err != nil {
			return response.FAIL
		}
		keySort := "post:sort:" + post.PostId
		conn.Redis.ZAdd(keySort, redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix())})
	}
	return response.SUCCESS
}

func GetPostDetail(detail *dao.PostDetail, stuId string) int {
	key := "post:id:" + detail.PostId
	stringDetail, err := conn.Redis.Get(key).Result()
	if err := json.Unmarshal([]byte(stringDetail), &detail.Post); err != nil {
		return response.FAIL
	}
	if err != nil {
		if err := conn.MySQL.First(&detail.Post).Error; err != nil {
			return response.FAIL
		}
		keySort := "post:sort:" + detail.PostId
		conn.Redis.ZAdd(keySort, redis.Z{Member: detail.PostId, Score: float64(detail.CreatedAt.Unix())})
	}
	_, detail.Star, detail.IsStar = liked.GetTotalAndLiked(detail.PostId, stuId)
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	return response.SUCCESS
}
