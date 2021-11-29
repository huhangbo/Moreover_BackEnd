package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"github.com/go-redis/redis"
)

func PublishLike(liked dao.Liked) int {
	if err := conn.MySQL.Create(&liked).Error; err != nil {
		return response.FAIL
	}
	if exist := conn.Redis.Exists("liked:" + liked.Parent).Val(); exist == 1 {
		conn.Redis.ZAdd("liked:"+liked.Parent, redis.Z{Member: liked.Publisher, Score: float64(liked.CreatedAt.Unix())})
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

func GetLikeByPage(current, size int, parentId string) (int, []dao.UserInfoBasic, bool) {
	var isEnd bool
	ids := conn.Redis.ZRevRange("liked:"+parentId, int64((current-1)*size), int64(current*size)-1).Val()
	if len(ids) == 0 {
		if err := conn.MySQL.Model(dao.Liked{}).Select("publisher").Where("parent = ?", parentId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&ids).Error; err != nil {
			return response.FAIL, nil, isEnd
		}
	}
	if len(ids) < size {
		isEnd = true
	}
	code, tmpBasic := GetKindDetail(ids)
	return code, tmpBasic, isEnd
}
