package user

import (
	"Moreover/internal/util"
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

func GetUserInfoDetail(stuId, publisher string) (int, model.UserInfoDetail) {
	var tmpUserDetail model.UserInfoDetail
	code, userInfo := GetUserInfo(stuId)
	if code != response.SUCCESS {
		return code, tmpUserDetail
	}
	_, totalFollower := util.GetTotalById(stuId, "follow", "follower")
	_, totalFan := util.GetTotalById(stuId, "follow", "fan")
	tmpUserDetail.UserInfo = userInfo
	tmpUserDetail.Follower = totalFollower
	tmpUserDetail.Fan = totalFan
	tmpUserDetail.IsFollow = util.IsPublished(stuId, "follow", "follower", "fan", publisher)
	return response.SUCCESS, tmpUserDetail
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

func GetKindDetail(likes []string) (int, []model.UserBasicInfo) {
	var likesDetail []model.UserBasicInfo
	for _, item := range likes {
		code, tmp := GetUserInfo(item)
		if code != response.SUCCESS {
			return code, likesDetail
		}
		likesDetail = append(likesDetail, tmp.UserBasicInfo)
	}
	return response.SUCCESS, likesDetail
}
