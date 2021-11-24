package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/util"
)

func PublishLike(liked dao.Liked) int {
	if err := conn.MySQL.Create(&liked).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(liked.Publisher, float64(liked.CreatedAt.Unix()), "liked:sort:"+liked.Parent) {
		return response.FAIL
	}
	return response.SUCCESS
}

func UnLike(liked dao.Liked) int {
	if err := conn.MySQL.Delete(&dao.Liked{}, "parent = ? AND publisher = ?", liked.Parent, liked.Publisher).Error; err != nil {
		return response.PasswordError
	}
	if _, err := conn.Redis.ZRem("liked:sort:"+liked.Parent, liked.Publisher).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

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
	code, tmpBasic = GetKindDetail(likes)
	return code, tmpBasic, tmpPage
}
