package liked

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
	"fmt"
	goRedis "github.com/go-redis/redis"
	"time"
)

const timeLikedExpiration = time.Hour * 24 * 7

func PublishLike(like model.Like) int {
	code := publishLikeToMysql(like)
	if code != response.SUCCESS {
		code = deleteLikeFromMysql(like.LikeId, 0)
		if code != response.SUCCESS {
			return code
		}
	}
	code = publishLikeToRedis(like)
	return code
}

func publishLikeToRedis(like model.Like) int {
	publishTime, _ := time.ParseInLocation("2006/01/02 15:05:06", like.UpdateTime, time.Local)
	sortComment := goRedis.Z{
		Score:  float64(publishTime.Unix()),
		Member: like.LikeId,
	}
	key := "liked:id:" + like.LikeId
	sortKey := "liked:sort:" + like.ParentId
	likeJson, err := json.Marshal(like)
	if err != nil {
		return response.ERROR
	}
	pipe := redis.DB.Pipeline()
	pipe.Set(key, likeJson, timeLikedExpiration)
	pipe.ZAdd(sortKey, sortComment)
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}

func publishLikeToMysql(like model.Like) int {
	sql := `INSERT INTO liked (create_time, update_time, like_id, parent_id, like_user, like_publisher, deleted)
			VALUES (:create_time, :update_time, :like_id, :parent_id, :like_user, :like_publisher, :deleted);`
	if _, err := mysql.DB.NamedExec(sql, like); err != nil {
		fmt.Println(err)
		return response.ERROR
	}
	return response.SUCCESS
}
