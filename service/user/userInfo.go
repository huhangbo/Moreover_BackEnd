package user

import (
	"Moreover/connent"
	"Moreover/dao"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"encoding/json"
	"time"
)

type DetailInfo struct {
	dao.UserInfo
	Follower int
	Fan      int
	IsFollow bool
}

const InfoExpiration = time.Hour * 24 * 7

func GetUserInfo(info *dao.UserInfo) int {
	code := getUserInfoFromRedis(info)
	if code != response.SUCCESS {
		connent.MySQL.First(info, &info.StudentId)
	}
	code = publishUserInfoToRedis(*info)
	return code
}

func GetUserInfoBasic(basic *dao.UserInfoBasic) int {
	var userInfo = dao.UserInfo{
		StudentId: basic.StudentId,
	}
	code := GetUserInfo(&userInfo)
	if code != response.SUCCESS {
		return code
	}
	*basic = dao.UserInfoBasic{
		StudentId: userInfo.StudentId,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Sex:       userInfo.Sex,
	}
	return response.SUCCESS
}

func GetUserInfoDetail(detail *DetailInfo, stuId string) int {
	tmpUserInfo := dao.UserInfo{
		StudentId: detail.StudentId,
	}
	GetUserInfo(&tmpUserInfo)
	(*detail).UserInfo = tmpUserInfo
	_, (*detail).Follower = util.GetTotalById(stuId, "follow", "follower")
	_, (*detail).Fan = util.GetTotalById(stuId, "follow", "fan")
	(*detail).IsFollow = util.IsPublished(stuId, "follow", "follower", "fan", detail.StudentId)
	return response.SUCCESS
}

func GetKindDetail(likes []string) (int, []dao.UserInfoBasic) {
	var likesDetail []dao.UserInfoBasic
	for _, item := range likes {
		var tmpUserInfo dao.UserInfoBasic
		tmpUserInfo.StudentId = item
		code := GetUserInfoBasic(&tmpUserInfo)
		if code != response.SUCCESS {
			return code, likesDetail
		}
		likesDetail = append(likesDetail, tmpUserInfo)
	}
	return response.SUCCESS, likesDetail
}

func UpdateUserInfo(info dao.UserInfo) int {
	if err := connent.MySQL.Model(&info).Updates(info).Error; err != nil {
		return response.FAIL
	}
	connent.MySQL.First(&info)
	publishUserInfoToRedis(info)
	return response.SUCCESS
}

func getUserInfoFromRedis(info *dao.UserInfo) int {
	key := "user:id:" + info.StudentId
	userInfoString, err := redis.DB.Get(key).Result()
	if err != nil {
		return response.FAIL
	}
	if err := json.Unmarshal([]byte(userInfoString), info); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func publishUserInfoToRedis(info dao.UserInfo) int {
	key := "user:id:" + info.StudentId
	tmpInfo, err := json.Marshal(info)
	if err != nil {
		return response.ERROR
	}
	if _, err := redis.DB.Set(key, string(tmpInfo), InfoExpiration).Result(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
