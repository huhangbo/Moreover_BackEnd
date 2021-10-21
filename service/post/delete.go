package post

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
)

func DeletePost(post dao.Post, stuId string) int {
	tmpPost := dao.Post{PostId: post.PostId}
	GetPost(&tmpPost)
	if tmpPost.Publisher != stuId {
		return response.AuthError
	}
	if err := conn.MySQL.Delete(&post).Error; err != nil {
		return response.FAIL
	}
	key := "post:id:" + post.PostId
	keySort := "post:sort:"
	keyTop := "post:sort:top"
	if !util.DeleteSortRedis(post.PostId, key, keySort, keyTop) {
		return response.FAIL
	}
	return response.SUCCESS
}
