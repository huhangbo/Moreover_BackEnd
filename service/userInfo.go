package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/util"
	"encoding/json"
	"time"
)

const InfoExpiration = time.Hour * 24 * 7

func GetUserInfo(info *dao.UserInfo) int {
	code := getUserInfoFromRedis(info)
	if code != response.SUCCESS {
		if err := conn.MySQL.Model(dao.UserInfo{}).Where("student_id = ?", info.StudentId).First(info).Error; err != nil {
			return response.UserNotExist
		}
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

func GetUserInfoBasicFollow(basic *dao.UserInfoBasicFollow, stuId string) int {
	var userInfo = dao.UserInfo{
		StudentId: basic.StudentId,
	}
	code := GetUserInfo(&userInfo)
	if code != response.SUCCESS {
		return code
	}
	(*basic).UserInfoBasic = dao.UserInfoBasic{
		StudentId: userInfo.StudentId,
		Nickname:  userInfo.Nickname,
		Avatar:    userInfo.Avatar,
		Sex:       userInfo.Sex,
	}
	_, _, (*basic).IsFollow = util.GetTotalAndIs("parent", (*basic).StudentId, "parent", stuId)
	return response.SUCCESS
}

func GetUserInfoDetail(detail *dao.UserInfoDetail, stuId string) int {
	tmpUserInfo := dao.UserInfo{
		StudentId: detail.StudentId,
	}
	if code := GetUserInfo(&tmpUserInfo); code != response.SUCCESS {
		return code
	}
	(*detail).UserInfo = tmpUserInfo
	_, (*detail).Follower = util.GetTotalById("publisher", (*detail).StudentId, "publisher")
	_, (*detail).Fan, (*detail).IsFollow = util.GetTotalAndIs("parent", (*detail).StudentId, "parent", stuId)
	return response.SUCCESS
}

func GetKindDetail(ids []string) (int, []dao.UserInfoBasic) {
	var likesDetail []dao.UserInfoBasic
	for _, item := range ids {
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
	if err := conn.MySQL.Model(&info).Updates(info).Error; err != nil {
		return response.FAIL
	}
	conn.MySQL.First(&info, info.StudentId)
	publishUserInfoToRedis(info)
	return response.SUCCESS
}

func getUserInfoFromRedis(info *dao.UserInfo) int {
	key := "user:id:" + info.StudentId
	userInfoString, err := conn.Redis.Get(key).Result()
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
	if _, err := conn.Redis.Set(key, string(tmpInfo), InfoExpiration).Result(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
