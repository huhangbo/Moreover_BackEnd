package user

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
)

func UpdateUserAvatar(avatar, userId string) int {
	code, oldUserInfo := GetUserInfo(userId)
	if code != response.SUCCESS {
		return code
	}
	oldUserInfo.Avatar = avatar
	return updateUserInfo(oldUserInfo)
}

func UpdateUserSex(sex, userId string) int {
	code, oldUserInfo := GetUserInfo(userId)
	if code != response.SUCCESS {
		return code
	}
	oldUserInfo.Sex = sex
	return updateUserInfo(oldUserInfo)
}

func UpdateUserNickname(nickname, userId string) int {
	code, oldUserInfo := GetUserInfo(userId)
	if code != response.SUCCESS {
		return code
	}
	oldUserInfo.Nickname = nickname
	return updateUserInfo(oldUserInfo)
}

func UpdateUserDescription(description, userId string) int {
	code, oldUserInfo := GetUserInfo(userId)
	if code != response.SUCCESS {
		return code
	}
	oldUserInfo.Description = description
	return updateUserInfo(oldUserInfo)
}

func updateUserInfo(info model.UserInfo) int {
	code := updateUserInfoToMysql(info)
	if code != response.SUCCESS {
		return code
	}
	return updateUserInfoToRedis(info)
}

func updateUserInfoToRedis(info model.UserInfo) int {
	key := "user:id:" + info.StudentID
	tmpInfo, err := json.Marshal(info)
	if err != nil {
		return response.ERROR
	}
	if _, err := redis.DB.Set(key, string(tmpInfo), InfoExpiration).Result(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func updateUserInfoToMysql(info model.UserInfo) int {
	sql := `UPDATE user_info
		    SET nickname = :nickname, sex = :sex, description = :description, avatar = :avatar, tag = :tag
			WHERE student_id = :student_id`
	if _, err := mysql.DB.NamedExec(sql, info); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
