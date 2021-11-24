package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/util"
	"github.com/go-redis/redis"
	"time"
)

const timeFollowExpiration = time.Hour * 24 * 7

func PublishFollow(follow dao.Follow) int {
	tmpUser := dao.User{StudentId: follow.Parent}
	if !IsUserExist(tmpUser) {
		return response.UserNotExist
	}
	if err := conn.MySQL.Create(&follow).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(follow.Parent, float64(follow.CreatedAt.Unix()), "publisher:sort:"+follow.Publisher) || !util.PublishSortRedis(follow.Publisher, float64(follow.CreatedAt.Unix()), "parent:sort:"+follow.Parent) {
		return response.FAIL
	}
	return response.SUCCESS
}

func Unfollow(follow dao.Follow) int {
	if err := conn.MySQL.Where("publisher = ? AND parent = ?", follow.Publisher, follow.Parent).Delete(&dao.Follow{}).Error; err != nil {
		return response.FAIL
	}
	keyFollow := "publisher:sort:" + follow.Publisher
	keyFan := "parent:sort:" + follow.Parent
	pipe := conn.Redis.Pipeline()
	pipe.ZRem(keyFollow, follow.Parent)
	pipe.ZRem(keyFan, follow.Publisher)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
func GetFollowById(current, size int, follower, followType, tmp string) (int, []dao.UserInfoBasic, model.Page) {
	code, total := util.GetTotalById(followType, follower, followType)
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
	code, tmpBasic = GetKindDetail(follows)
	return code, tmpBasic, tmpPage
}

func GetFollowByIdFromRedis(current, size int, follower, category string) (int, []string) {
	code, fans := util.GetIdsByPageFromRedis(current, size, follower, category)
	if code != response.SUCCESS {
		return code, fans
	}
	return code, fans
}

func GetTotalFollow(follower string) (error, []string) {
	key := "parent:sort:" + follower
	followers, _ := conn.Redis.ZRange(key, 0, -1).Result()
	if len(followers) == 0 {
		if err := conn.MySQL.Model(&dao.Follow{}).Select("publisher").Find(&followers).Error; err != nil {
			return err, followers
		}
	}
	return nil, followers
}
func SyncFollowToRedis(follower, category, tmp string) {
	var follows []dao.Follow
	if err := conn.MySQL.Model(&dao.Follow{}).Select(tmp, "created_at").Where(category+" = ?", follower).Find(&follows).Error; err != nil {
		return
	}
	key := category + ":sort:" + follower
	pipe := conn.Redis.Pipeline()
	if category == "publisher" {
		for _, item := range follows {
			pipe.ZAdd(key, redis.Z{Member: item.Parent, Score: float64(item.CreatedAt.Unix())})
		}
	} else {
		for _, item := range follows {
			pipe.ZAdd(key, redis.Z{Member: item.Publisher, Score: float64(item.CreatedAt.Unix())})
		}
	}
	pipe.Expire(key, timeFollowExpiration)
	if _, err := pipe.Exec(); err != nil {
	}
}
