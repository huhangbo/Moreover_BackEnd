package comment

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/service/util"
)

func SyncCommentSortRedis(parentId string) {
	var comments []dao.Comment
	if err := conn.MySQL.Where("parent_id = ?", parentId).Find(&comments).Error; err != nil {
		return
	}
	for _, item := range comments {
		sortKey := "comment:sort:" + item.ParentId
		util.PublishSortRedis(item.CommentId, float64(item.UpdatedAt.Unix()), sortKey)
	}
}
