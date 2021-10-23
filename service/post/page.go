package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"github.com/go-redis/redis"
)

func GetPostByPage(current, size int, stuId string) (int, []dao.PostDetail, model.Page) {
	var posts []dao.PostDetail
	sortKey := "post:sort:"
	total, _ := conn.Redis.ZCard(sortKey).Result()
	var tmpPage model.Page
	if total == 0 {
		var tmpPosts []dao.Post
		if err := conn.MySQL.Model(&dao.Post{}).Find(&tmpPosts).Order("created_at DESC").Error; err != nil {
			return response.ERROR, posts, tmpPage
		}
		total = int64(len(tmpPosts))
		pipe := conn.Redis.Pipeline()
		for _, item := range tmpPosts {
			pipe.ZAdd(sortKey, redis.Z{Member: item.PostId, Score: float64(item.CreatedAt.Unix())})
		}
		if _, err := pipe.Exec(); err != nil {
			return response.ERROR, posts, tmpPage
		}
	}
	tmpPage = model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if int(total) < (current-1)*size {
		return response.ERROR, posts, tmpPage
	}
	_, postIds := util.GetIdsByPageFromRedis(current, size, "", "post")
	for _, item := range postIds {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		GetPostDetail(&tmpDetail, stuId)
		posts = append(posts, tmpDetail)
	}
	return response.SUCCESS, posts, tmpPage
}

func GetPostByPublisher(current, size int, stuId string) (int, []dao.Post, model.Page) {
	var posts []dao.Post
	var postIds []string
	var total int64
	if err := conn.MySQL.Model(&dao.Post{}).Where("publisher = ?", stuId).Count(&total).Error; err != nil {
		return response.FAIL, posts, model.Page{}
	}
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	conn.MySQL.Model(&dao.Post{}).Select("post_id").Where("publisher = ?", stuId).Limit(size).Offset((current - 1) * size).Order("created_at desc").Find(&postIds)
	for i := 0; i < len(postIds); i++ {
		tmpPost := dao.Post{PostId: postIds[i]}
		GetPost(&tmpPost)
		posts = append(posts, tmpPost)
	}
	return response.SUCCESS, posts, tmpPage
}
func GetPostByTop(current, size int, stuId string) (int, []dao.PostDetail, model.Page) {
	var posts []dao.PostDetail
	sortKey := "post:sort:top"
	total, _ := conn.Redis.ZCard(sortKey).Result()
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if int(total) < (current-1)*size {
		return response.ERROR, posts, tmpPage
	}
	_, postIds := util.GetIdsByPageFromRedis(current, size, "top", "post")
	for _, item := range postIds {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		GetPostDetail(&tmpDetail, stuId)
		posts = append(posts, tmpDetail)
	}
	return response.SUCCESS, posts, tmpPage
}
