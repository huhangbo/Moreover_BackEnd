package liked

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
)

func DeleteLikeById(likeId, stuId string) int {
	code, tmpLike := GetLikeById(likeId)
	if code != response.SUCCESS {
		return code
	}
	if stuId != tmpLike.LikePublisher {
		return response.AuthError
	}
	code = deleteLikeFromMysql(likeId, 1)
	if code != response.SUCCESS {
		return code
	}
	return deleteLikeFromRedis(likeId, tmpLike.ParentId)
}

func deleteLikeFromRedis(likeId, parentId string) int {
	idKey := "liked:id:" + likeId
	sortParentKey := "liked:sort:" + parentId
	pipe := redis.DB.Pipeline()
	pipe.Del(idKey)
	pipe.ZRem(sortParentKey, likeId)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func deleteLikeFromMysql(likeId string, deleted int) int {
	sql := `UPDATE liked
			SET deleted = ?
			WHERE like_id = ?`

	if _, err := mysql.DB.Exec(sql, deleted, likeId); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
