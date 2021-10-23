package util

import (
	"Moreover/conn"
	"Moreover/pkg/response"
)

func GetTotalById(parentId, kind, parent string) (int, int) {
	sortKey := kind + ":sort:" + parentId
	total, err := conn.Redis.ZCard(sortKey).Result()
	if kind == "publisher" || kind == "parent" {
		kind = "follow"
	}
	if err != nil || total == 0 {
		if err := conn.MySQL.Table(kind).Where(parent+" = ?", parentId).Count(&total).Error; err != nil {
			return response.FAIL, int(total)
		}
		if total != 0 {
			return response.NotFound, int(total)
		}
	}
	return response.SUCCESS, int(total)
}

func GetIdsByPageFromRedis(current, size int, parentId, kind string) (int, []string) {
	sortKey := kind + ":sort:" + parentId
	ids, err := conn.Redis.ZRevRange(sortKey, int64((current-1)*size), int64(current*size)).Result()
	if err != nil {
		return response.ERROR, ids
	}
	return response.SUCCESS, ids
}
