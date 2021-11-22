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

const (
	postExpiration = time.Hour * 24 * 7
	sortKey        = "post:sort:"
)

func PublishPost(post dao.Post) int {
	post.Picture = util.ArrayToString(post.Pictures)
	if err := conn.MySQL.Create(&post).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	postJson, _ := json.Marshal(post)
	pipe := conn.Redis.Pipeline()
	pipe.ZAdd(sortKey, redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix())})
	pipe.ZAdd("post:sort:top", redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix()) / 10000})
	pipe.Set(key, string(postJson), postExpiration)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
