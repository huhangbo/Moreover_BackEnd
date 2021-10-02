package comment

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"fmt"
)

func DeleteCommentById(commentId, stuId string) int {
	code, tmpComment := GetCommentById(commentId)
	if code != response.SUCCESS {
		return code
	}
	if stuId != tmpComment.Publisher {
		return response.AuthError
	}
	code = deleteCommentFromMysql(commentId)
	if code != response.SUCCESS {
		return code
	}
	return deleteCommentFromRedis(commentId, tmpComment.ParentID)
}

func deleteCommentFromRedis(commentId, parentId string) int {
	idKey := "comment:id:" + commentId
	sortParentKey := "comment:sort:" + parentId
	pipe := redis.DB.Pipeline()
	pipe.Del(idKey)
	pipe.ZRem(sortParentKey, commentId)
	if _, err := pipe.Exec(); err != nil {
		fmt.Printf("delete activity from redis fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}

func deleteCommentFromMysql(commentId string) int {
	sql := `UPDATE comment
			SET deleted = 1
			WHERE comment_id = ?
			OR parent_id = ?`

	if _, err := mysql.DB.Exec(sql, commentId, commentId); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
