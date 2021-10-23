package follow

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
)

func GetFollowById(current, size int, follower, followType, tmp string) (int, []dao.UserInfoBasic, model.Page) {
	code, total := util.GetTotalById(follower, followType, followType)
	var tmpBasic []dao.UserInfoBasic
	var follows []string
	tmpPage := model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: total/size + 1,
	}
	if code != response.SUCCESS {
		if code != response.NotFound {
			return code, tmpBasic, tmpPage
		}
		SyncFollowToRedis(follower, followType, tmp)
	}
	if (current-1)*size > total {
		return code, tmpBasic, tmpPage
	}
	code, follows = GetFollowByIdFromRedis(current, size, follower, followType)
	if code != response.SUCCESS || len(follows) == 0 {
		if err := conn.MySQL.Model(&dao.Follow{}).Select(tmp).Where(followType+" = ?", follower).Limit(size).Offset((current - 1) * size).Find(&follows).Order("created_at DESC").Error; err != nil {
			return code, tmpBasic, tmpPage
		}
	}
	code, tmpBasic = user.GetKindDetail(follows)
	return code, tmpBasic, tmpPage
}

func GetFollowByIdFromRedis(current, size int, follower, category string) (int, []string) {
	code, fans := util.GetIdsByPageFromRedis(current, size, follower, category)
	if code != response.SUCCESS {
		return code, fans
	}
	return code, fans
}
