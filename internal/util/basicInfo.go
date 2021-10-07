package util

import (
	"Moreover/internal/pkg/user"
	"Moreover/model"
	"Moreover/pkg/response"
)

func GetKindDetail(likes []string) (int, []model.UserBasicInfo) {
	var likesDetail []model.UserBasicInfo
	for _, item := range likes {
		code, tmp := user.GetUserInfo(item)
		if code != response.SUCCESS {
			return code, likesDetail
		}
		likesDetail = append(likesDetail, tmp.UserBasicInfo)
	}
	return response.SUCCESS, likesDetail
}
