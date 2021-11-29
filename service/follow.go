package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"github.com/go-redis/redis"
	"time"
)

const (
	timeFollowExpiration = time.Hour * 24 * 7
)

func PublishFollow(follow dao.Follow) int {
	tmpUser := dao.User{
		StudentId: follow.Parent,
	}
	if !IsUserExist(tmpUser) {
		return response.UserNotExist
	}
	if err := conn.MySQL.Create(&follow).Error; err != nil {
		return response.FAIL
	}
	if exist1 := conn.Redis.Exists("publisher:" + follow.Publisher).Val(); exist1 == 1 {
		conn.Redis.ZAdd("publisher:"+follow.Publisher, redis.Z{Member: follow.Parent, Score: float64(follow.CreatedAt.Unix())})
	}
	if exist2 := conn.Redis.Exists("parent:" + follow.Parent).Val(); exist2 == 1 {
		conn.Redis.ZAdd("parent:"+follow.Parent, redis.Z{Member: follow.Publisher, Score: float64(follow.CreatedAt.Unix())})
	}
	return response.SUCCESS
}

func Unfollow(follow dao.Follow) int {
	if err := conn.MySQL.Where("publisher = ? AND parent = ?", follow.Publisher, follow.Parent).Delete(&dao.Follow{}).Error; err != nil {
		return response.FAIL
	}
	keyFollow := "publisher:" + follow.Publisher
	keyFan := "parent:" + follow.Parent
	pipe := conn.Redis.Pipeline()
	pipe.ZRem(keyFollow, follow.Parent)
	pipe.ZRem(keyFan, follow.Publisher)
	_, _ = pipe.Exec()
	return response.SUCCESS
}

func GetFollowById(current, size int, follower, followType, tmp string) (int, []dao.UserInfoBasic, bool) {
	var (
		isEnd bool
	)
	ids, _ := conn.Redis.ZRange(followType+":"+follower, int64((current-1)*size), int64(current*size-1)).Result()
	if len(ids) == 0 {
		wg.Add(1)
		go SyncFollowToRedis(follower, followType, tmp)
		if err := conn.MySQL.Model(&dao.Follow{}).Select(tmp).Where(followType+" = ?", follower).Limit(size).Offset((current - 1) * size).Find(&ids).Order("created_at DESC").Error; err != nil {
			return response.FAIL, nil, isEnd
		}
	}
	if len(ids) < size {
		isEnd = true
	}
	code, tmpBasic := GetKindDetail(ids)
	wg.Wait()
	return code, tmpBasic, isEnd
}

func GetTotalFollow(follower string) (error, []string) {
	key := "parent:" + follower
	followers, _ := conn.Redis.ZRange(key, 0, -1).Result()
	if len(followers) == 0 {
		if err := conn.MySQL.Model(&dao.Follow{}).Select("publisher").Find(&followers).Error; err != nil {
			return err, followers
		}
	}
	return nil, followers
}

func SyncFollowToRedis(follower, category, tmp string) {
	defer wg.Done()
	var follows []dao.Follow
	if err := conn.MySQL.Model(&dao.Follow{}).Select(tmp, "created_at").Where(category+" = ?", follower).Find(&follows).Error; err != nil {
		return
	}
	key := category + ":" + follower
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
	_, _ = pipe.Exec()
}
