package comment

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	goRedis "github.com/go-redis/redis"
	"time"

	"fmt"
)

const commentExpiration = time.Hour * 24 * 7

func PublishComment(comment model.Comment) int {
	code := publishCommentToMysql(comment)
	if code != response.SUCCESS {
		return code
	}
	return publishCommentToRedis(comment)
}

func publishCommentToRedis(comment model.Comment) int {
	jsonActivity, err := json.Marshal(comment)
	publishTime, _ := time.ParseInLocation("2006/01/02 15:05:06", comment.CreateTime, time.Local)
	if err != nil {
		return response.ERROR
	}
	sortComment := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: comment.CommentId,
	}
	idKey := "comment:id:" + comment.CommentId
	sortKey := "comment:sort:" + comment.ParentID
	pipe := redis.DB.Pipeline()
	pipe.ZAdd(sortKey, sortComment)
	pipe.Set(idKey, string(jsonActivity), commentExpiration)
	if _, err := pipe.Exec(); err != nil {
		fmt.Printf("insert comment to redis fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}

func publishCommentToMysql(comment model.Comment) int {
	sql := `INSERT INTO comment (create_time, update_time, comment_id, publisher, replier, parent_id, message)
			VALUES (:create_time, :update_time, :comment_id, :publisher, :replier, :parent_id, :message)`
	if _, err := mysql.DB.NamedExec(sql, comment); err != nil {
		fmt.Printf("insert comment to mysql fail, err: %v\n", err)
		return response.ERROR
	}
	return response.SUCCESS
}
