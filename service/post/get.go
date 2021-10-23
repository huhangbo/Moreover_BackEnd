package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"Moreover/service/util"
	"encoding/json"
)

func GetPost(post *dao.Post) int {
	key := "post:id:" + post.PostId
	postString, err := conn.Redis.Get(key).Result()
	_ = json.Unmarshal([]byte(postString), post)
	if err != nil {
		if err := conn.MySQL.First(post).Error; err != nil {
			return response.FAIL
		}
		post.Pictures = util.StringToArray(post.Picture)
		postJson, _ := json.Marshal(post)
		if _, err := conn.Redis.Set(key, string(postJson), postExpiration).Result(); err != nil {
			return response.FAIL
		}
	}
	return response.SUCCESS
}

func GetPostDetail(detail *dao.PostDetail, stuId string) int {
	if code := GetPost(&detail.Post); code != response.SUCCESS {
		return code
	}
	_, detail.Star, detail.IsStar = liked.GetTotalAndLiked(detail.PostId, stuId)
	_, detail.Comments = util.GetTotalById(detail.PostId, "comment", "parent_id")
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	return response.SUCCESS
}
