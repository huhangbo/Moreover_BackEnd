package comment

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	goRedis "github.com/go-redis/redis"
)

func GetCommentById(commentId string) (int, model.Comment) {
	code, comment := getCommentByIdFromRedis(commentId)
	if code != response.SUCCESS {
		code, comment = getCommentByIdFromMysql(commentId)
		if code != response.SUCCESS {
			return code, comment
		}
		code = publishCommentToRedis(comment)
	}
	return code, comment
}

func GetCommentsByIdPage(current, size int, commentId string) (int, []model.Comment, model.Page) {
	_, totalComment := GetTotalCommentById(commentId)
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
	code, ids := getCommentIdsByPageFromRedis(current, size, commentId)
	if code != response.SUCCESS {
		code, comments = getCommentsByPageFromMysql(current, size, commentId)
		return code, comments, tmpPage
	}
	for i := 0; i < len(ids); i++ {
		code, tmpComment := GetCommentById(ids[i])
		if code != response.SUCCESS {
			return code, comments, tmpPage
		}
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

func GetTotalCommentById(parentId string) (int, int) {
	code, total := getTotalByIdFromRedis(parentId)
	if code != response.SUCCESS {
		code, total = getTotalByIdFromMysql(parentId)
	}
	return code, total
}

func getTotalByIdFromRedis(parentId string) (int, int) {
	sortKey := "comment:sort:" + parentId
	total, err := redis.DB.ZCard(sortKey).Result()
	if err != nil {
		return response.ERROR, int(total)
	}
	return response.SUCCESS, int(total)
}

func getTotalByIdFromMysql(parentId string) (int, int) {
	var totalComment int
	sql := `SELECT COUNT(*)
			FROM comment
			WHERE parent_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&totalComment, sql, parentId); err != nil {
		return response.ERROR, totalComment
	}
	return response.SUCCESS, totalComment
}

func getCommentIdsByPageFromRedis(current, size int, parentId string) (int, []string) {
	sortKey := "comment:sort:" + parentId
	rangeOpt := goRedis.ZRangeBy{
		Min:    "-",
		Max:    "+",
		Offset: int64((current - 1) * size),
		Count:  int64(size),
	}
	ids, err := redis.DB.ZRangeByLex(sortKey, rangeOpt).Result()
	if err != nil {
		return response.ERROR, ids
	}
	return response.SUCCESS, ids
}

func getCommentsByPageFromMysql(current, size int, parentId string) (int, []model.Comment) {
	var comments []model.Comment
	sql := `SELECT * FROM comment
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY publish_time
			LIMIT ?, ?`
	if err := mysql.DB.Get(comments, sql, parentId, (current-1)*size, size); err != nil {
		return response.ERROR, comments
	}
	return response.SUCCESS, comments
}
