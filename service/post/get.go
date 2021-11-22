package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/follow"
	"Moreover/service/user"
	"Moreover/service/util"
	"encoding/json"
)

func GetPost(post *dao.Post) int {
	key := "post:id:" + post.PostId
	postString, err := conn.Redis.Get(key).Result()
	_ = json.Unmarshal([]byte(postString), post)
	if err != nil {
		if err := conn.MySQL.Model(dao.Post{}).Where("post_id = ?", post.PostId).First(post).Error; err != nil {
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
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.PostId, "parent", stuId)
	_, detail.Comments = util.GetTotalById("comment", detail.PostId, "parent_id")
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	return response.SUCCESS
}

func GetFollowPost(current, size int, stuId string) (int, bool, []dao.PostDetail) {
	var (
		posts   []dao.PostDetail
		postIds []string
		isEnd   bool
	)
	err, followers := follow.GetTotalFollow(stuId)
	if err != nil {
		return response.ParamError, isEnd, posts
	}
	if err := conn.MySQL.Model(dao.Post{}).Where("publisher IN ?", followers).Select("post_id").Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&postIds).Error; err != nil {
		return response.FAIL, isEnd, posts
	}
	if len(postIds) < size {
		isEnd = true
	}
	for i := 0; i < len(postIds); i++ {
		tmpPost := dao.PostDetail{Post: dao.Post{PostId: postIds[i]}}
		if code := GetPostDetail(&tmpPost, stuId); code != response.SUCCESS {
			return response.FAIL, isEnd, posts
		}
		posts = append(posts, tmpPost)
	}
	return response.SUCCESS, isEnd, posts
}
