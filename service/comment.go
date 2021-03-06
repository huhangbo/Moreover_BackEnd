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
	commentExpiration = time.Hour * 24 * 7
	commentSortKey    = "comment:sort:"
	commentIdKey      = "comment:id:"
)

func PublishComment(comment dao.Comment) int {
	if err := conn.MySQL.Create(&comment).Error; err != nil {
		return response.FAIL
	}
	comment.PublishedAt = comment.CreatedAt.Unix()
	conn.Redis.ZAdd(commentSortKey+comment.ParentId, redis.Z{Score: float64(comment.PublishedAt), Member: comment.CommentId})
	publishCommentToRedis(comment)
	return response.SUCCESS
}

func publishCommentToRedis(comment dao.Comment) int {
	jsonActivity, _ := json.Marshal(comment)
	key := commentIdKey + comment.CommentId
	if err := conn.Redis.Set(key, string(jsonActivity), commentExpiration); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func DeleteComment(comment dao.Comment, stuId string) int {
	code := GetCommentById(&comment)
	if code != response.SUCCESS {
		return code
	}
	if comment.Publisher != stuId {
		return response.AuthError
	}
	if err := conn.MySQL.Where("parent_id = ?", comment.CommentId).Or("comment_id = ?", comment.CommentId).Delete(dao.Comment{}).Error; err != nil {
		return response.ParamError
	}
	return deleteCommentFromRedis(comment)
}

func deleteCommentFromRedis(comment dao.Comment) int {
	idKey := commentIdKey + comment.CommentId
	sortKey := commentSortKey + comment.CommentId
	sortParentKey := commentSortKey + comment.ParentId
	pipe := conn.Redis.Pipeline()
	pipe.Del(idKey)
	pipe.Del(sortKey)
	pipe.ZRem(sortParentKey, comment.CommentId)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
func GetCommentById(comment *dao.Comment) int {
	code := getCommentByIdFromRedis(comment)
	if code != response.SUCCESS {
		if err := conn.MySQL.Where("comment_id = ?", comment.CommentId).First(comment).Error; err != nil {
			return response.FAIL
		}
		comment.PublishedAt = comment.CreatedAt.Unix()
		publishCommentToRedis(*comment)
		return response.SUCCESS
	}
	return code
}

func getCommentByIdFromRedis(comment *dao.Comment) int {
	commentString, err := conn.Redis.Get(commentIdKey + comment.CommentId).Result()
	if err != nil || commentString == "" {
		return response.FAIL
	}
	if err := json.Unmarshal([]byte(commentString), comment); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func GetCommentByIdPage(current, size int, parentId, stuId string) (int, []dao.CommentDetail, bool) {
	var (
		commentsDetail []dao.CommentDetail
		isEnd          bool
	)
	ids := conn.Redis.ZRevRange(commentSortKey+parentId, int64((current-1)*size), int64(current*size-1)).Val()
	if len(ids) == 0 {
		wg.Add(1)
		go SyncCommentSortRedis(parentId)
		if err := conn.MySQL.Model(&dao.Comment{}).Select("comment_id").Where("parent_id = ?", parentId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&ids).Error; err != nil {
			return response.FAIL, nil, isEnd
		}
	}
	for i := 0; i < len(ids); i++ {
		tmpCommentDetail := dao.CommentDetail{
			Comment: dao.Comment{
				CommentId: ids[i],
			},
		}
		if code := GetCommentById(&(tmpCommentDetail.Comment)); code != response.SUCCESS {
			return code, nil, isEnd
		}
		tmpCommentDetail.PublisherInfo = dao.UserInfoBasic{
			StudentId: tmpCommentDetail.Publisher,
		}
		GetUserInfoBasic(&(tmpCommentDetail.PublisherInfo))
		_, tmpCommentDetail.Star, tmpCommentDetail.IsStar = util.GetTotalAndIs("liked", tmpCommentDetail.CommentId, "parent_id", stuId)
		commentsDetail = append(commentsDetail, tmpCommentDetail)
	}
	if len(ids) < size {
		isEnd = true
	}
	wg.Wait()
	return response.SUCCESS, commentsDetail, isEnd
}

func GetCommentChildrenByPage(current, size int, commentId, stuId string) (int, []dao.CommentChild, bool) {
	var commentChildren []dao.CommentChild
	code, childrenDetail, isEnd := GetCommentByIdPage(current, size, commentId, stuId)
	for _, item := range childrenDetail {
		tmpCommentChild := dao.CommentChild{
			CommentDetail: item,
			ReplierInfo: dao.UserInfoBasic{
				StudentId: item.Replier,
			},
		}
		GetUserInfoBasic(&(tmpCommentChild.ReplierInfo))
		commentChildren = append(commentChildren, tmpCommentChild)
	}
	return code, commentChildren, isEnd
}

func SyncCommentSortRedis(parentId string) {
	defer wg.Done()
	var (
		tmpIds []struct {
			CommentId string
			CreatedAt time.Time
		}
		tmpZs []redis.Z
	)
	if err := conn.MySQL.Model(&dao.Comment{}).Select("comment_id AND created_at").Where("parent_id = ?", parentId).Find(&tmpIds).Error; err != nil {
		return
	}
	for _, item := range tmpIds {
		tmpZs = append(tmpZs, redis.Z{Member: item, Score: float64(item.CreatedAt.Unix())})
	}
	conn.Redis.ZAdd(commentSortKey+parentId, tmpZs...)
}
