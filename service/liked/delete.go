package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func UnLike(liked dao.Liked) int {
	if err := conn.MySQL.Delete(&dao.Liked{}, "parent = ? AND publisher = ?", liked.Parent, liked.Publisher).Error; err != nil {
		return response.PasswordError
	}
	if _, err := conn.Redis.ZRem("liked:sort:"+liked.Parent, liked.Publisher).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
