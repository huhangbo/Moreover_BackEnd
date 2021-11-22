package comment

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"encoding/json"
)

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
