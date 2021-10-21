package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func UnLike(liked dao.Liked) int {
	if err := conn.MySQL.Unscoped().Delete(&liked).Error; err != nil {
		return response.FAIL
	}
	if _, err := conn.Redis.ZRem("like:sort:"+liked.ParentId, liked.Publisher).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
