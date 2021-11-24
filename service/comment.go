package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/util"
	"encoding/json"
	"time"
)

const commentExpiration = time.Hour * 24 * 7

func PublishComment(comment dao.Comment) int {
	if err := conn.MySQL.Create(&comment).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(comment.CommentId, float64(comment.CreatedAt.Unix()), "comment:sort:"+comment.ParentId) {
		return response.FAIL
	}
	publishCommentToRedis(comment)
	return response.SUCCESS
}

func publishCommentToRedis(comment dao.Comment) int {
	jsonActivity, _ := json.Marshal(comment)
	key := "comment:id:" + comment.CommentId
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
	idKey := "comment:id:" + comment.CommentId
	sortKey := "comment:sort:" + comment.CommentId
	sortParentKey := "comment:sort:" + comment.ParentId
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
		publishCommentToRedis(*comment)
		return response.SUCCESS
	}
	return code
}

func getCommentByIdFromRedis(comment *dao.Comment) int {
	commentString, err := conn.Redis.Get("comment:id:" + comment.CommentId).Result()
	if err != nil || commentString == "" {
		return response.FAIL
	}
	if err := json.Unmarshal([]byte(commentString), comment); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
func GetCommentByIdPage(current, size int, parentId, stuId string) (int, []dao.CommentDetail, model.Page) {
	code, total := util.GetTotalById("comment", parentId, "parent_id")
	if code == response.NotFound {
		SyncCommentSortRedis(parentId)
	}
	var commentsDetail []dao.CommentDetail
	var commentIds []string
	var tmpPage = model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: (total / size) + 1,
	}
	if total == 0 {
		return response.SUCCESS, commentsDetail, tmpPage
	}
	if (current-1)*size > total {
		return response.ParamError, commentsDetail, tmpPage
	}
	code, commentIds = util.GetIdsByPageFromRedis(current, size, parentId, "comment")
	if code != response.SUCCESS || (len(commentIds) == 0 && code == 200) {
		conn.MySQL.Model(&dao.Comment{}).Select("comment_id").Where("parent_id = ?", parentId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&commentIds)
	}
	for i := 0; i < len(commentIds); i++ {
		tmpCommentDetail := dao.CommentDetail{
			Comment: dao.Comment{
				CommentId: commentIds[i],
			},
		}
		code = GetCommentById(&(tmpCommentDetail.Comment))
		if code != response.SUCCESS {
			return code, commentsDetail, tmpPage
		}
		tmpCommentDetail.PublisherInfo = dao.UserInfoBasic{
			StudentId: tmpCommentDetail.Publisher,
		}
		GetUserInfoBasic(&(tmpCommentDetail.PublisherInfo))
		_, tmpCommentDetail.Star, tmpCommentDetail.IsStart = util.GetTotalAndIs("liked", tmpCommentDetail.CommentId, "parent_id", stuId)
		commentsDetail = append(commentsDetail, tmpCommentDetail)
	}
	return response.SUCCESS, commentsDetail, tmpPage
}

func GetCommentChildrenByPage(current, size int, commentId, stuId string) (int, []dao.CommentChild, model.Page) {
	var commentChildren []dao.CommentChild
	code, childrenDetail, tmpPage := GetCommentByIdPage(current, size, commentId, stuId)
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
	return code, commentChildren, tmpPage
}

func SyncCommentSortRedis(parentId string) {
	var comments []dao.Comment
	if err := conn.MySQL.Where("parent_id = ?", parentId).Find(&comments).Error; err != nil {
		return
	}
	for _, item := range comments {
		sortKey := "comment:sort:" + item.ParentId
		util.PublishSortRedis(item.CommentId, float64(item.UpdatedAt.Unix()), sortKey)
	}
}
