package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
)

func DeletePost(post dao.Post) int {
	tmpPost := dao.Post{PostId: post.PostId}
	if code := GetPost(&tmpPost); code != response.SUCCESS {
		return code
	}
	if err := conn.MySQL.Where("post_id = ? AND publisher = ?", post.PostId, post.Publisher).Delete(&dao.Activity{}).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	keyTop := "post:sort:top"
	if !util.DeleteSortRedis(post.PostId, key, sortKey, keyTop) {
		return response.FAIL
	}
	return response.SUCCESS
}
