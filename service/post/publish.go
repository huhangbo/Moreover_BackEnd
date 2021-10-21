package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

const talkExpiration = time.Hour * 24 * 7

func PublishPost(post dao.Post) int {
	post.Picture = util.ArrayToString(post.Pictures)
	if err := conn.MySQL.Create(&post).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	talkJson, _ := json.Marshal(post)
	conn.Redis.Set(key, string(talkJson), talkExpiration)
	return response.SUCCESS
}

func TopPost(post dao.PostDetail) int {
	key := "post:sort:"
	score := util.GetTopScore(post.Star, post.CreatedAt)
	if err := conn.Redis.ZAdd(key, redis.Z{Member: post.PostId, Score: score}); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
