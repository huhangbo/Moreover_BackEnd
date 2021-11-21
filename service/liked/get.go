package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
)

func GetLikeByPage(current, size int, parentId string) (int, []dao.UserInfoBasic, model.Page) {
	code, total := util.GetTotalById("liked", parentId, "parent")
	var tmpBasic []dao.UserInfoBasic
	var likes []string
	tmpPage := model.Page{Current: current, PageSize: size, Total: total, TotalPage: total/size + 1}
	if code != response.SUCCESS {
		return code, tmpBasic, tmpPage
	}
	if (current-1)*size > total {
		return code, tmpBasic, tmpPage
	}
	code, likes = util.GetIdsByPageFromRedis(current, size, parentId, "liked")
	if code != response.SUCCESS {
		if err := conn.MySQL.Model(dao.Liked{}).Select("publisher").Where("parent = ?", parentId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&likes).Error; err != nil {
			return response.FAIL, tmpBasic, tmpPage
		}
	}
	code, tmpBasic = user.GetKindDetail(likes)
	return code, tmpBasic, tmpPage
}
