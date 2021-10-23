package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"time"
)

const timeLikedExpiration = time.Hour * 24 * 7

func PublishLike(liked dao.Liked) int {
	if err := conn.MySQL.Create(&liked).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(liked.Publisher, float64(liked.CreatedAt.Unix()), "liked:sort:"+liked.ParentId) {
		return response.FAIL
	}
	return response.SUCCESS
}
