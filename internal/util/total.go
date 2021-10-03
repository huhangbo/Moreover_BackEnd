package util

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	goRedis "github.com/go-redis/redis"
)

func GetTotalById(parentId, category string) (int, int) {
	code, total := getTotalByIdFromRedis(parentId, category)
	if code != response.SUCCESS {
		code, total = getTotalByIdFromMysql(parentId, category)
	}
	return code, total
}

func GetIdsByPageFromRedis(current, size int, parentId, category string) (int, []string) {
	sortKey := category + ":sort:" + parentId
	rangeOpt := goRedis.ZRangeBy{
		Min:    "-",
		Max:    "+",
		Offset: int64((current - 1) * size),
		Count:  int64(size),
	}
	ids, err := redis.DB.ZRangeByLex(sortKey, rangeOpt).Result()
	if err != nil {
		return response.ERROR, ids
	}
	return response.SUCCESS, ids
}

func getTotalByIdFromRedis(parentId, category string) (int, int) {
	sortKey := category + ":sort:" + parentId
	total, err := redis.DB.ZCard(sortKey).Result()
	if err != nil || total == 0 {
		return response.ERROR, int(total)
	}
	return response.SUCCESS, int(total)
}

func getTotalByIdFromMysql(parentId, category string) (int, int) {
	var total int
	sql := `SELECT COUNT(*)
			FROM ` + category + `
			WHERE parent_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&total, sql, parentId); err != nil {
		return response.ERROR, total
	}
	return response.SUCCESS, total
}
