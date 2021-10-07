package liked

import (
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
	"time"
)

func SyncLikeMysqlToRedis(parentId string) int {
	var tmpLikes []model.Like
	key := "like:sort:" + parentId
	sql := `SELECT parent_id, like_publisher, update_time
			FROM liked
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY update_time DESC`
	if err := mysql.DB.Select(&tmpLikes, sql, parentId); err != nil {
		return response.ERROR
	}
	pipe := redis.DB.Pipeline()
	for _, item := range tmpLikes {
		tmpPublishTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.UpdateTime, time.Local)
		pipe.ZAdd(key, goRedis.Z{Member: item.LikePublisher, Score: float64(tmpPublishTime.Unix())})
	}
	if _, err := pipe.Exec(); err != nil {
		return response.ERROR
	}
	return response.SUCCESS
}
