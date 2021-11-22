package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"github.com/go-redis/redis"
	"time"
)

func GetPostByPage(current, size int, stuId string) (int, []dao.PostDetail, model.Page) {
	var (
		posts   []dao.PostDetail
		tmpPage model.Page
	)
	total, _ := conn.Redis.ZCard(sortKey).Result()
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
		pipe.ZAdd(sortKey, tmpZs...)
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
