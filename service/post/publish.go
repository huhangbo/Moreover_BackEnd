package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
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
