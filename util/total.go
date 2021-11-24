package util

import (
	"Moreover/conn"
	"Moreover/pkg/response"
	"github.com/go-redis/redis"
	"time"
)

type SortSet struct {
	Publisher string
	CreatedAt time.Time
}

func GetTotalById(kind, parentId, parent string) (int, int) {
	sortKey := kind + ":sort:" + parentId
	total, err := conn.Redis.ZCard(sortKey).Result()
	if kind == "publisher" || kind == "parent" {
		kind = "follow"
	}
	if err != nil || total == 0 {
		if err := conn.MySQL.Table(kind).Where(parent+" = ? AND deleted_at = ?", parentId, 0).Count(&total).Error; err != nil {
			return response.FAIL, int(total)
		}
	}
	return response.SUCCESS, int(total)
}

func GetIdsByPageFromRedis(current, size int, parentId, kind string) (int, []string) {
	sortKey := kind + ":sort:" + parentId
	ids, err := conn.Redis.ZRevRange(sortKey, int64((current-1)*size), int64(current*size)-1).Result()
	if err != nil {
		return response.NotFound, ids
	}
	return response.SUCCESS, ids
}

func GetTotalAndIs(kind, parentId, parent, publisher string) (int, int, bool) {
	var (
		sortKey    = kind + ":sort:" + parentId
		tmpSorts   []SortSet
		is         bool
		total, err = conn.Redis.ZCard(sortKey).Result()
	)
	if err != nil || total == 0 {
		if kind == "publisher" || kind == "parent" {
			kind = "follow"
		}
		if err := conn.MySQL.Table(kind).Select("publisher", "created_at").Where(parent+" = ? AND deleted_at = ?", parentId, 0).Find(&tmpSorts).Error; err != nil {
			return response.FAIL, len(tmpSorts), is
		}
		if len(tmpSorts) != 0 {
			var tmpZs []redis.Z
			for _, item := range tmpSorts {
				tmpZ := redis.Z{Member: item.Publisher, Score: float64(item.CreatedAt.Unix())}
				tmpZs = append(tmpZs, tmpZ)
				if item.Publisher == publisher {
					is = true
				}
			}
			pipe := conn.Redis.Pipeline()
			pipe.ZAdd(sortKey, tmpZs...)
			pipe.Expire(sortKey, time.Hour*7*24)
			if _, err := pipe.Exec(); err != nil {
				return response.FAIL, len(tmpSorts), is
			}
		}
		return response.SUCCESS, len(tmpSorts), is
	}
	count, _ := conn.Redis.ZScore(sortKey, publisher).Result()
	if count != 0 {
		is = true
	}
	return response.SUCCESS, int(total), is
}
