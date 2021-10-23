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

const postExpiration = time.Hour * 24 * 7

func PublishPost(post dao.Post) int {
	post.Picture = util.ArrayToString(post.Pictures)
	if err := conn.MySQL.Create(&post).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	sortKey := "post:sort:"
	postJson, _ := json.Marshal(post)
	conn.Redis.ZAdd(sortKey, redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix())})
	conn.Redis.Set(key, string(postJson), postExpiration)
	return response.SUCCESS
}
