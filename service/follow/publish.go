package follow

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"time"
)

const timeFollowExpiration = time.Hour * 24 * 7

func PublishFollow(follow dao.Follow) int {
	if err := conn.MySQL.Create(&follow).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(follow.Parent, float64(time.Now().Unix()), "publisher:sort:"+follow.Publisher) || !util.PublishSortRedis(follow.Publisher, float64(time.Now().Unix()), "parent:sort:"+follow.Parent) {
		return response.FAIL
	}
	return response.SUCCESS
}