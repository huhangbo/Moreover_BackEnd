package util

import (
	"Moreover/conn"
)

func IsPublished(parentId, kind, parent, publishType, publisher string) bool {
	if isPublishedFromRedis(parentId, kind, publisher) == true {
		return true
	}
	return isPublishedFromMysql(parentId, kind, parent, publishType, publisher)
}

func isPublishedFromRedis(parentId, kind, publisher string) bool {
	sortKey := kind + ":sort:" + parentId
	total, _ := conn.Redis.ZScore(sortKey, publisher).Result()
	if total > 0 {
		return true
	}
	return false
}

func isPublishedFromMysql(parentId, kind, parent, publishType, publisher string) bool {
	var tmpNum int64
	if err := conn.MySQL.Table(kind).Where(parent+" = ?", parentId).Where(publishType+"= ? ", publisher).Count(&tmpNum).Error; err != nil || tmpNum == 0 {
		return false
	}
	return true
}
