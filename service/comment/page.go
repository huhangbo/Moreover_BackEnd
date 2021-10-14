package comment

import (
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
)

func GetCommentByIdPage(current, size int, commentId, stuId string) (int, []model.CommentDetail, model.Page) {
	_, totalComment := util.GetTotalById(commentId, "comment", "parent_id")
	var comments []model.CommentDetail
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
				publishCommentToRedis(comments[i].Comment)
				tmpUserBasic := dao.UserInfoBasic{
					StudentId: comments[i].Publisher,
				}
				user.GetUserInfoBasic(&tmpUserBasic)
				comments[i].PublisherInfo = tmpUserBasic
			}
		}
		return code, comments, tmpPage
	}
	for i := 0; i < len(ids); i++ {
		var tmpComment model.CommentDetail
		code, tmpComment.Comment = GetCommentById(ids[i])
		if code != response.SUCCESS {
			return code, comments, tmpPage
		}
		tmpUserBasic := dao.UserInfoBasic{
			StudentId: tmpComment.Publisher,
		}
		user.GetUserInfoBasic(&tmpUserBasic)
		tmpComment.PublisherInfo = tmpUserBasic
		_, tmpComment.Star = util.GetTotalById(commentId, "likes", "parent_id")
		tmpComment.IsStart = util.IsPublished(commentId, "likes", "parent_id", "like_publisher", stuId)
		comments = append(comments, tmpComment)
	}
	return response.SUCCESS, comments, tmpPage
}

func GetPreChildCById(size int, commentId, stuId string) (int, model.ChildComment) {
	var code int
	var children model.ChildComment
	var page model.Page
	code, children.Comments, page = GetCommentByIdPage(1, size, commentId, stuId)
	children.Total = page.Total
	if code != response.SUCCESS {
		return code, children
	}
	return code, children
}

func getCommentsByPageFromMysql(current, size int, parentId string) (int, []model.CommentDetail) {
	var comments []model.CommentDetail
	sql := `SELECT * FROM comment
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY update_time
			LIMIT ?, ?`
	if err := mysql.DB.Select(&comments, sql, parentId, (current-1)*size, size); err != nil {
		return response.ERROR, comments
	}
	return response.SUCCESS, comments
}
