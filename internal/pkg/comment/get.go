package comment

import (
	"Moreover/internal/pkg/user"
	"Moreover/internal/util"
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
			return code, comment
		}
	}
	return code, comment
}

func GetCommentsByIdPage(current, size int, commentId string) (int, []model.Comment, model.Page) {
	_, totalComment := util.GetTotalById(commentId, "comment")
	var comments []model.Comment
	var tmpPage = model.Page{
		Current:   current,
		PageSize:  size,
		Total:     totalComment,
		TotalPage: (totalComment / size) + 1,
	}
	if totalComment == 0 {
		return response.SUCCESS, comments, tmpPage
	}
	if (current-1)*size > totalComment {
		return response.ParamError, comments, tmpPage
	}
	code, ids := util.GetIdsByPageFromRedis(current, size, commentId, "comment")
	if code != response.SUCCESS || (len(ids) == 0 && code == 200) {
		code, comments = getCommentsByPageFromMysql(current, size, commentId)
		if code == response.SUCCESS {
			for i := 0; i < len(comments); i++ {
				publishCommentToRedis(comments[i])
				_, comments[i].PublisherInfo = user.GetUserInfo(comments[i].Publisher)
				comments[i].PublisherInfo.Description = ""
			}
		}
		return code, comments, tmpPage
	}
	for i := 0; i < len(ids); i++ {
		code, tmpComment := GetCommentById(ids[i])
		if code != response.SUCCESS {
			return code, comments, tmpPage
		}
		_, tmpComment.PublisherInfo = user.GetUserInfo(tmpComment.Publisher)
		tmpComment.PublisherInfo.Description = ""
		comments = append(comments, tmpComment)
	}
	return response.SUCCESS, comments, tmpPage
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

func getCommentsByPageFromMysql(current, size int, parentId string) (int, []model.Comment) {
	var comments []model.Comment
	sql := `SELECT * FROM comment
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY publish_time
			LIMIT ?, ?`
	if err := mysql.DB.Select(&comments, sql, parentId, (current-1)*size, size); err != nil {
		return response.ERROR, comments
	}
	return response.SUCCESS, comments
}
