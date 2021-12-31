package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

const (
	postExpiration = time.Hour * 24 * 7
	sortPostKey    = "post:sort:"
)

func PublishPost(post dao.Post) int {
	post.Picture = util.ArrayToString(post.Pictures)
	if err := conn.MySQL.Create(&post).Error; err != nil {
		return response.FAIL
	}
	post.PublishedAt = post.CreatedAt.Unix()
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
	post.Pictures = util.StringToArray(post.Picture)
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
	if err := conn.MySQL.Where("post_id = ? AND publisher = ?", post.PostId, post.Publisher).Delete(&dao.Post{}).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	keyTop := "post:sort:top"
	if !dao.DeleteSortRedis(post.PostId, key, sortPostKey, keyTop) {
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
		post.Pictures = util.StringToArray(post.Picture)
		post.PublishedAt = post.CreatedAt.Unix()
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
	GetUserInfoBasic(&(detail.PublisherInfo))
	return response.SUCCESS
}

func GetFollowPost(current, size int, stuId string) (int, []dao.PostDetail, bool) {
	var (
		posts   []dao.PostDetail
		postIds []string
		isEnd   bool
	)
	err, followers := GetTotalFollow(stuId)
	if err != nil {
		return response.ParamError, nil, isEnd
	}
	if err := conn.MySQL.Model(dao.Post{}).Where("publisher IN ?", followers).Select("post_id").Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&postIds).Error; err != nil {
		return response.FAIL, nil, isEnd
	}
	if len(postIds) < size {
		isEnd = true
	}
	for i := 0; i < len(postIds); i++ {
		tmpPost := dao.PostDetail{Post: dao.Post{PostId: postIds[i]}}
		if code := GetPostDetail(&tmpPost, stuId); code != response.SUCCESS {
			return response.FAIL, nil, isEnd
		}
		posts = append(posts, tmpPost)
	}
	return response.SUCCESS, posts, isEnd
}

func GetPostByPage(current, size int, stuId string) (int, []dao.PostDetail, bool) {
	var (
		posts []dao.PostDetail
		isEnd bool
	)
	ids := conn.Redis.ZRevRange("post:sort:", int64((current-1)*size), int64(current*size)-1).Val()
	if len(ids) == 0 && current == 1 {
		wg.Add(1)
		go syncPostToRedis()
		if err := conn.MySQL.Model(&dao.Post{}).Select("post_id").Find(&ids).Limit(size).Offset((current - 1) * size).Error; err != nil {
			return response.FAIL, nil, isEnd
		}
	}
	if len(ids) < size {
		isEnd = true
	}
	for _, item := range ids {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		GetPostDetail(&tmpDetail, stuId)
		posts = append(posts, tmpDetail)
	}
	wg.Wait()
	return response.SUCCESS, posts, isEnd
}

func GetPostByPublisher(current, size int, stuId, userId string) (int, []dao.PostDetail, bool) {
	var (
		posts []dao.PostDetail
		ids   []string
		isEnd bool
	)
	if err := conn.MySQL.Model(&dao.Post{}).Select("post_id").Where("publisher = ?", userId).Limit(size).Offset((current - 1) * size).Order("created_at desc").Find(&ids).Error; err != nil {
		return response.FAIL, nil, isEnd
	}
	if len(ids) < size {
		isEnd = true
	}
	for i := 0; i < len(ids); i++ {
		tmpPost := dao.PostDetail{Post: dao.Post{PostId: ids[i]}}
		if code := GetPostDetail(&tmpPost, stuId); code != response.SUCCESS {
			return response.FAIL, nil, isEnd
		}
		posts = append(posts, tmpPost)
	}
	return response.SUCCESS, posts, isEnd
}

func GetPostByTop(current, size int, stuId string) (int, []dao.PostDetail, bool) {
	var (
		posts []dao.PostDetail
		isEnd bool
	)
	ids := conn.Redis.ZRevRange("post:top:", int64((current-1)*size), int64(current*size)-1).Val()
	for _, item := range ids {
		tmpDetail := dao.PostDetail{Post: dao.Post{PostId: item}}
		if code := GetPostDetail(&tmpDetail, stuId); code != response.SUCCESS {
			return code, nil, isEnd
		}
		posts = append(posts, tmpDetail)
	}
	if len(ids) < size {
		isEnd = true
	}
	return response.SUCCESS, posts, isEnd
}

func syncPostToRedis() {
	defer wg.Done()
	var (
		tmpIds []struct {
			PostId    string
			CreatedAt time.Time
		}
		tmpZs []redis.Z
	)
	if err := conn.MySQL.Model(&dao.Post{}).Find(&tmpIds).Error; err != nil {
		return
	}
	for i := 0; i < len(tmpIds); i++ {
		tmpZs = append(tmpZs, redis.Z{Member: tmpIds[i], Score: float64(tmpIds[i].CreatedAt.Unix())})
	}
	conn.Redis.ZAdd(sortPostKey, tmpZs...)
}
