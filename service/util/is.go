package util

import (
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
)

func IsPublished(parentId, kind, parent, publishType, publisher string) bool {
	if isPublishedFromRedis(parentId, kind, publisher) == true {
		return true
	}
	return isPublishedFromMysql(parentId, kind, parent, publishType, publisher)
}

func isPublishedFromRedis(parentId, kind, publisher string) bool {
	sortKey := kind + ":sort:" + parentId
	total, _ := redis.DB.ZScore(sortKey, publisher).Result()
	if total > 0 {
		return true
	}
	return false
}

func isPublishedFromMysql(parentId, kind, parent, publishType, publisher string) bool {
	var tmpNum int
	sql := `SELECT COUNT(*)
			FROM ` + kind + `
			WHERE ` + parent + ` = ?
			AND ` + publishType + ` = ?
			AND deleted = 0`
	if err := mysql.DB.Get(tmpNum, sql, parentId, publisher); err != nil || tmpNum == 0 {
		return false
	}
	return true
}
