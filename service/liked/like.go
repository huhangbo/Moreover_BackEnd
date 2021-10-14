package liked

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
	"time"
)

const timeLikedExpiration = time.Hour * 24 * 7

func PublishLike(like model.Like) int {
	code := publishLikeToRedis(like)
	if code != response.SUCCESS {
		return code
	}
	code = publishLikeToMysql(like)
	if code != response.SUCCESS {
		UnLikeFromRedis(like.ParentId, like.LikePublisher)
	}
	return code
}

func publishLikeToRedis(like model.Like) int {
	publishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", like.UpdateTime, time.Local)
	sortKey := "liked:sort:" + like.ParentId
	sortComment := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: like.LikePublisher,
	}
	pipe := redis.DB.Pipeline()
	pipe.ZAdd(sortKey, sortComment)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func publishLikeToMysql(like model.Like) int {
	sql := `INSERT INTO liked (create_time, update_time, parent_id, like_user, like_publisher)
			VALUES (:create_time, :update_time,:parent_id, :like_user, :like_publisher);`
	if _, err := mysql.DB.NamedExec(sql, like); err != nil {
		code := UnLikeFromMysql(like.ParentId, like.LikePublisher, 0)
		return code
	}
	return response.SUCCESS
}
