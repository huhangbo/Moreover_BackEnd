package comment

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
	"time"
)

const commentExpiration = time.Hour * 24 * 7

func PublishComment(comment dao.Comment) int {
	if err := conn.MySQL.Create(&comment).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(comment.CommentId, float64(time.Now().Unix()), "comment:sort:"+comment.ParentId) {
		return response.FAIL
	}
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
