package follow

import (
	"Moreover/internal/util"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/response"
)

func GetFollowById(current, size int, follower, followType string) (int, []model.UserBasicInfo, model.Page) {
	code, total := util.GetTotalById(follower, followType)
	var tmpBasic []model.UserBasicInfo
	var follows []string
	tmpPage := model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: total/size + 1,
	}
	if code != response.SUCCESS || (current-1)*size > total {
		return code, tmpBasic, tmpPage
	}
	code, follows = GetFollowByIdFromRedis(current, size, follower, followType)
	if code != response.SUCCESS || len(follows) == 0 {
		code, follows = GetFollowByIdFromMysql(current, size, follower, followType)
		if code == response.SUCCESS {
			SyncFollowMysqlToRedis(follower, followType)
		}
		code, tmpBasic = util.GetKindDetail(follows)
		return code, tmpBasic, tmpPage
	}
	code, tmpBasic = util.GetKindDetail(follows)
	return code, tmpBasic, tmpPage
}

func GetFollowByIdFromRedis(current, size int, follower, category string) (int, []string) {
	code, fans := util.GetIdsByPageFromRedis(current, size, follower, category)
	if code != response.SUCCESS {
		return code, fans
	}
	return code, fans
}

func GetFollowByIdFromMysql(current, size int, follower, category string) (int, []string) {
	tmp := category
	if category == "follower" {
		category = "fan"
	} else {
		category = "follower"
	}
	var follows []string
	sql := `SELECT ` + tmp +
		` FROM follow
			WHERE ` + category + ` = ?
			AND deleted = 0
			ORDER BY update_time DESC 
			LIMIT ?, ?`
	if err := mysql.DB.Select(&follows, sql, follower, (current-1)*size, size); err != nil {
		return response.ERROR, follows
	}
	return response.SUCCESS, follows
}
