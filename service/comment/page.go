package comment

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"Moreover/service/util"
)

func GetCommentByIdPage(current, size int, commentId, stuId string) (int, []dao.CommentDetail, model.Page) {
	code, total := util.GetTotalById(commentId, "comment", "parent_id")
	if code == response.NotFound {
		SyncCommentSortRedis(commentId)
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
	code, commentIds = util.GetIdsByPageFromRedis(current, size, commentId, "comment")
	if code != response.SUCCESS || (len(commentIds) == 0 && code == 200) {
		conn.MySQL.Model(&dao.Comment{}).Select("comment_id").Where("parent_id = ?", commentId).Limit(size).Offset((current - 1) * size).Find(&commentIds)
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
		user.GetUserInfoBasic(&(tmpCommentDetail.PublisherInfo))
		_, tmpCommentDetail.Star, tmpCommentDetail.IsStart = liked.GetTotalAndLiked(tmpCommentDetail.CommentId, stuId)
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
		user.GetUserInfoBasic(&(tmpCommentChild.ReplierInfo))
		commentChildren = append(commentChildren, tmpCommentChild)
	}
	return code, commentChildren, tmpPage
}
