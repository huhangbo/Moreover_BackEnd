package liked

import (
	"Moreover/internal/util"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
)

func GetLikeById(current, size int, parentId string) (int, []model.UserBasicInfo, model.Page) {
	code, total := util.GetTotalById(parentId, "liked")
	var tmpBasic []model.UserBasicInfo
	var likes []string
	tmpPage := model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: total/size + 1,
	}
	if code != response.SUCCESS || (current-1)*size > total {
		return code, tmpBasic, tmpPage
	}
	code, likes = util.GetIdsByPageFromRedis(current, size, parentId, "liked")
	if code != response.SUCCESS || len(likes) == 0 {
		code, likes = getLikeByIdFromMysql(current, size, parentId)
		if code == response.SUCCESS {
			SyncLikeMysqlToRedis(parentId)
		}
		code, tmpBasic = util.GetKindDetail(likes)
		return code, tmpBasic, tmpPage
	}
	code, tmpBasic = util.GetKindDetail(likes)
	return code, tmpBasic, tmpPage
}

func getLikeByIdFromMysql(current, size int, parentId string) (int, []string) {
	var likes []string
	sql := `SELECT like_publisher
			FROM liked
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY update_time DESC 
			LIMIT ?, ?`
	if err := mysql.DB.Get(&likes, sql, parentId, (current-1)*size, size); err != nil {
		return response.ERROR, likes
	}
	return response.SUCCESS, likes
}
