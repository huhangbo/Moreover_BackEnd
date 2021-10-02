package user

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"time"
)

const InfoExpiration = time.Hour * 24 * 7

const (
	DefaultAvatar      = "https://moreover-1305054989.cos.ap-nanjing.myqcloud.com/author.jpg"
	DefaultSex         = "未知"
	DefaultName        = "取个名字吧"
	DefaultDescription = "添加一句话描述下自己吧"
)

func GetUserInfo(stuId string) (int, model.UserInfo) {
	code, userInfo := getUserInfoFromRedis(stuId)
	if code != response.SUCCESS {
		code, userInfo = getUserInfoFromMysql(stuId)
	}
	code = updateUserInfoToRedis(userInfo)
	return code, userInfo
}

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

func getUserInfoFromRedis(stuId string) (int, model.UserInfo) {
	var userInfo model.UserInfo
	key := "user:id:" + stuId
	userInfoString, err := redis.DB.Get(key).Result()
	if err != nil {
		return response.ERROR, userInfo
	}
	if err := json.Unmarshal([]byte(userInfoString), &userInfoString); err != nil {
		return response.ERROR, userInfo
	}
	return response.SUCCESS, userInfo
}

func getUserInfoFromMysql(stuId string) (int, model.UserInfo) {
	var userInfo model.UserInfo
	sql := `SELECT student_id, nickname, sex, description, avatar, tag
			FROM user_info
			WHERE student_id = ?`
	if err := mysql.DB.Get(&userInfo, sql, stuId); err != nil {
		return response.ERROR, userInfo
	}
	return response.SUCCESS, userInfo
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
