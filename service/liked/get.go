package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
	"github.com/go-redis/redis"
)

func GetLikeByPage(current, size int, parentId string) (int, []dao.UserInfoBasic, model.Page) {
	code, total := util.GetTotalById(parentId, "liked", "parent_id")
	var tmpBasic []dao.UserInfoBasic
	var likes []string
	tmpPage := model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: total/size + 1,
	}
	if code != response.SUCCESS {
		return code, tmpBasic, tmpPage
	}
	if (current-1)*size > total {
		return code, tmpBasic, tmpPage
	}
	code, likes = util.GetIdsByPageFromRedis(current, size, parentId, "liked")
	if code != response.SUCCESS || len(likes) == 0 {
		if err := conn.MySQL.Model(dao.Liked{}).Select("publisher").Where("parent_id = ?", parentId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&likes).Error; err != nil {
			return response.FAIL, tmpBasic, tmpPage
		}
	}
	code, tmpBasic = user.GetKindDetail(likes)
	return code, tmpBasic, tmpPage
}

func GetTotalAndLiked(parentId, publisher string) (int, int, bool) {
	sortKey := "liked:sort:" + parentId
	var isLiked bool
	total, err := conn.Redis.ZCard(sortKey).Result()
	if err != nil || total == 0 {
		var likes []dao.Liked
		if err := conn.MySQL.Model(dao.Liked{}).Where("parent_id = ?", parentId).Find(&likes).Error; err != nil {
			return response.FAIL, len(likes), isLiked
		}
		if len(likes) != 0 {
			SyncLikeToRedis(likes)
			pipe := conn.Redis.Pipeline()
			for _, item := range likes {
				pipe.ZAdd(sortKey, redis.Z{Member: item.Publisher, Score: float64(item.CreatedAt.Unix())})
				if item.Publisher == publisher {
					isLiked = true
				}
			}
			pipe.Expire(sortKey, timeLikedExpiration)
			if _, err := pipe.Exec(); err != nil {
				return response.FAIL, len(likes), isLiked
			}
		}
		return response.SUCCESS, len(likes), false
	}
	count, _ := conn.Redis.ZScore(sortKey, publisher).Result()
	if count > 0 {
		isLiked = true
	}
	return response.SUCCESS, int(total), isLiked
}
