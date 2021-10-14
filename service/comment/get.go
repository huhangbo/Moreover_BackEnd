package comment

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
)

func GetCommentById(commentId string) (int, model.Comment) {
	code, comment := getCommentByIdFromRedis(commentId)
	if code != response.SUCCESS {
		code, comment = getCommentByIdFromMysql(commentId)
		if code == response.SUCCESS {
			publishCommentToRedis(comment)
		}
	}
	return code, comment
}

func getCommentByIdFromRedis(commentId string) (int, model.Comment) {
	var comment model.Comment
	idKey := "comment:id:" + commentId
	commentString, err := redis.DB.Get(idKey).Result()
	if err != nil {
		return response.ParamError, comment
	}
	if err := json.Unmarshal([]byte(commentString), &comment); err != nil {
		return response.ERROR, comment
	}
	return response.SUCCESS, comment
}

func getCommentByIdFromMysql(commentId string) (int, model.Comment) {
	var comment model.Comment
	sql := `SELECT * FROM comment
			WHERE comment_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&comment, sql, commentId); err != nil {
		fmt.Printf("get activity by id from mysql fail, err: %v\n", err)
		return response.ERROR, comment
	}
	return response.SUCCESS, comment
}
