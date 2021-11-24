package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	util2 "Moreover/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

const (
	postExpiration = time.Hour * 24 * 7
	sortPostKey    = "post:sort:"
)

func PublishPost(post dao.Post) int {
	post.Picture = util2.ArrayToString(post.Pictures)
	if err := conn.MySQL.Create(&post).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	postJson, _ := json.Marshal(post)
	pipe := conn.Redis.Pipeline()
	pipe.ZAdd(sortPostKey, redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix())})
	pipe.ZAdd("post:sort:top", redis.Z{Member: post.PostId, Score: float64(post.CreatedAt.Unix()) / 10000})
	pipe.Set(key, string(postJson), postExpiration)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func UpdatePost(post dao.Post) int {
	if err := conn.MySQL.Model(&dao.Post{}).Where("post_id = ? AND publisher = ?", post.PostId, post.Publisher).Updates(post).Error; err != nil {
		return response.ERROR
	}
	if err := conn.MySQL.Model(&dao.Post{}).Where("post_id = ?", post.PostId).First(&post).Error; err != nil {
		return response.ERROR
	}
	post.Pictures = util2.StringToArray(post.Picture)
	key := "post:id:" + post.PostId
	postJson, _ := json.Marshal(post)
	if _, err := conn.Redis.Set(key, string(postJson), postExpiration).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func DeletePost(post dao.Post) int {
	tmpPost := dao.Post{PostId: post.PostId}
	if code := GetPost(&tmpPost); code != response.SUCCESS {
		return code
	}
	if err := conn.MySQL.Where("post_id = ? AND publisher = ?", post.PostId, post.Publisher).Delete(&dao.Activity{}).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	keyTop := "post:sort:top"
	if !util2.DeleteSortRedis(post.PostId, key, sortPostKey, keyTop) {
		return response.FAIL
	}
	return response.SUCCESS
}

func GetPost(post *dao.Post) int {
	key := "post:id:" + post.PostId
	postString, err := conn.Redis.Get(key).Result()
	_ = json.Unmarshal([]byte(postString), post)
	if err != nil {
		if err := conn.MySQL.Model(dao.Post{}).Where("post_id = ?", post.PostId).First(post).Error; err != nil {
			return response.FAIL
		}
		post.Pictures = util2.StringToArray(post.Picture)
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
	_, detail.Star, detail.IsStar = util2.GetTotalAndIs("liked", detail.PostId, "parent", stuId)
	_, detail.Comments = util2.GetTotalById("comment", detail.PostId, "parent_id")
	detail.PublisherInfo.StudentId = detail.Publisher
	GetUserInfoBasic(&(detail.PublisherInfo))
	return response.SUCCESS
}

func GetFollowPost(current, size int, stuId string) (int, bool, []dao.PostDetail) {
	var (
		posts   []dao.PostDetail
		postIds []string
		isEnd   bool
	)
	err, followers := GetTotalFollow(stuId)
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

func GetPostByPage(current, size int, stuId string) (int, []dao.PostDetail, model.Page) {
	var (
		posts   []dao.PostDetail
		tmpPage model.Page
	)
	total, _ := conn.Redis.ZCard(sortPostKey).Result()
	if total == 0 {
		var (
			tmpIds []struct {
				PostId    string
				CreatedAt time.Time
			}
			tmpZs []redis.Z
		)
		if err := conn.MySQL.Model(&dao.Post{}).Find(&tmpIds).Error; err != nil {
			return response.ERROR, posts, tmpPage
		}
		total = int64(len(tmpIds))
		pipe := conn.Redis.Pipeline()
		for i := 0; i < len(tmpIds); i++ {
			tmpZs = append(tmpZs, redis.Z{Member: tmpIds[i], Score: float64(tmpIds[i].CreatedAt.Unix())})
		}
		pipe.ZAdd(sortPostKey, tmpZs...)
		if _, err := pipe.Exec(); err != nil {
			return response.ERROR, posts, tmpPage
		}
	}
	tmpPage = model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if int(total) < (current-1)*size {
		return response.ERROR, posts, tmpPage
	}
	_, postIds := util2.GetIdsByPageFromRedis(current, size, "", "post")
	for _, item := range postIds {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		GetPostDetail(&tmpDetail, stuId)
		posts = append(posts, tmpDetail)
	}
	return response.SUCCESS, posts, tmpPage
}

func GetPostByPublisher(current, size int, stuId string) (int, []dao.Post, model.Page) {
	var (
		posts   []dao.Post
		postIds []string
		total   int64
	)
	if err := conn.MySQL.Model(&dao.Post{}).Where("publisher = ?", stuId).Count(&total).Error; err != nil {
		return response.FAIL, posts, model.Page{}
	}
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if err := conn.MySQL.Model(&dao.Post{}).Select("post_id").Where("publisher = ?", stuId).Limit(size).Offset((current - 1) * size).Order("created_at desc").Find(&postIds).Error; err != nil {
		return response.FAIL, posts, tmpPage
	}
	for i := 0; i < len(postIds); i++ {
		tmpPost := dao.Post{PostId: postIds[i]}
		GetPost(&tmpPost)
		posts = append(posts, tmpPost)
	}
	return response.SUCCESS, posts, tmpPage
}
func GetPostByTop(current, size int, stuId string) (int, []dao.PostDetail, model.Page) {
	var posts []dao.PostDetail
	total, _ := conn.Redis.ZCard(sortPostKey).Result()
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if int(total) < (current-1)*size {
		return response.ERROR, posts, tmpPage
	}
	_, postIds := util2.GetIdsByPageFromRedis(current, size, "top", "post")
	for _, item := range postIds {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		GetPostDetail(&tmpDetail, stuId)
		posts = append(posts, tmpDetail)
	}
	return response.SUCCESS, posts, tmpPage
}
