package follow

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func Unfollow(follow dao.Follow) int {
	if err := conn.MySQL.Unscoped().Delete(&follow).Error; err != nil {
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
