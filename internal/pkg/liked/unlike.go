package liked

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
)

func UnLike(parentId, stuId string) int {
	code := UnLikeFromRedis(parentId, stuId)
	if code != response.SUCCESS {
		return code
	}
	code = UnLikeFromMysql(parentId, stuId, 1)
	if code != response.SUCCESS {
		publishLikeToRedis(model.Like{
			ParentId:      parentId,
			LikePublisher: stuId,
		})
	}
	return code
}

func UnLikeFromRedis(parentId, stuId string) int {
	key := "like:sort:" + parentId
	pipe := redis.DB.Pipeline()
	pipe.ZRem(key, stuId)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func UnLikeFromMysql(parentId, stuId string, state int) int {
	sql := `UPDATE liked
			SET deleted = ?
			WHERE parent_id = ?
			AND like_publisher = ?`
	if _, err := mysql.DB.Exec(sql, state, parentId, stuId); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
