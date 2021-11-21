package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
)

func UpdatePost(post dao.Post) int {
	if err := conn.MySQL.Model(&dao.Post{}).Where("post_id = ? AND publisher = ?", post.PostId, post.Publisher).Updates(post).Error; err != nil {
		return response.ERROR
	}
	if err := conn.MySQL.Model(&dao.Post{}).Where("post_id = ?", post.PostId).First(&post).Error; err != nil {
		return response.ERROR
	}
	post.Pictures = util.StringToArray(post.Picture)
	key := "post:id:" + post.PostId
	postJson, _ := json.Marshal(post)
	if _, err := conn.Redis.Set(key, string(postJson), postExpiration).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
