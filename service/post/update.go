package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
)

func UpdatePost(post dao.Post, stuId string) int {
	tmpPost := dao.Post{PostId: post.PostId}
	if code := GetPost(&tmpPost); code != response.SUCCESS {
		return response.FAIL
	}
	if tmpPost.Publisher != stuId {
		return response.AuthError
	}
	if err := conn.MySQL.Model(&dao.Post{}).Where("post_id = ?", post.PostId).Updates(post).Error; err != nil {
		return response.ERROR
	}
	if err := conn.MySQL.First(&post).Error; err != nil {
		return response.ERROR
	}
	post.Pictures = util.StringToArray(post.Picture)
	key := "post:id:" + post.PostId
	postJson, _ := json.Marshal(post)
	conn.Redis.Set(key, string(postJson), postExpiration)
	return response.SUCCESS
}
