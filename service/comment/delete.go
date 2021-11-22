package comment

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

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
