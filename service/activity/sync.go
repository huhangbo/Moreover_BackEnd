package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/service/util"
)

func SyncActivitySortRedis() {
	var activities []dao.Activity
	if err := conn.MySQL.Find(&activities).Error; err != nil {
		return
	}
	for _, item := range activities {
		sortCategoryKey := "activity:sort:" + item.Category
		sortKey := "activity:sort:"
		util.PublishSortRedis(item.ActivityId, float64(item.UpdatedAt.Unix()), sortKey, sortCategoryKey)
	}
}
