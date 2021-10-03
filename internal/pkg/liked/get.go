package liked

import (
	"Moreover/internal/pkg/user"
	"Moreover/internal/util"
	"Moreover/model"
	"Moreover/pkg/mysql"
	"Moreover/pkg/redis"
	"Moreover/pkg/response"
	"encoding/json"
)

func GetLikeById(likeId string) (int, model.Like) {
	code, like := getLikeByIdFromRedis(likeId)
	if code != response.SUCCESS {
		code, like = getLikeByIdFromMysql(likeId)
		if code == response.SUCCESS {
			publishLikeToRedis(like)
			return code, like
		}
	}
	return code, like
}

func getLikeByIdFromRedis(likeId string) (int, model.Like) {
	var like model.Like
	key := "liked:id:" + likeId
	likeString, err := redis.DB.Get(key).Result()
	if err != nil {
		return response.ParamError, like
	}
	if err := json.Unmarshal([]byte(likeString), &like); err != nil {
		return response.ERROR, like
	}
	return response.SUCCESS, like
}

func getLikeByIdFromMysql(likeId string) (int, model.Like) {
	var like model.Like
	sql := `SELECT * FROM liked
			WHERE like_id = ?
			AND deleted = 0`
	if err := mysql.DB.Get(&like, sql, likeId); err != nil {
		return response.ERROR, like
	}
	return response.SUCCESS, like
}

func GetLikesByPage(current, size int, parentId string) (int, []model.Like, model.Page) {
	_, totalComment := util.GetTotalById(parentId, "liked")
	var likes []model.Like
	var tmpPage = model.Page{
		Current:   current,
		PageSize:  size,
		Total:     totalComment,
		TotalPage: (totalComment / size) + 1,
	}
	if totalComment == 0 {
		return response.SUCCESS, likes, tmpPage
	}
	if (current-1)*size > totalComment {
		return response.ParamError, likes, tmpPage
	}
	code, ids := util.GetIdsByPageFromRedis(current, size, parentId, "liked")
	if code != response.SUCCESS || (len(ids) == 0 && code == 200) {
		code, likes = getLikesByPageFromMysql(current, size, parentId)
		if code == response.SUCCESS {
			for i := 0; i < len(likes); i++ {
				publishLikeToRedis(likes[i])
				_, likes[i].PublisherInfo = user.GetUserInfo(likes[i].LikePublisher)
				likes[i].PublisherInfo.Description = ""
			}
		}
		return code, likes, tmpPage
	}
	for i := 0; i < len(ids); i++ {
		code, tmpLike := GetLikeById(ids[i])
		if code != response.SUCCESS {
			return code, likes, tmpPage
		}
		_, tmpLike.PublisherInfo = user.GetUserInfo(tmpLike.LikePublisher)
		tmpLike.PublisherInfo.Description = ""
		likes = append(likes, tmpLike)
	}
	return response.SUCCESS, likes, tmpPage
}

func getLikesByPageFromMysql(current, size int, parentId string) (int, []model.Like) {
	var likes []model.Like
	sql := `SELECT * FROM liked
			WHERE parent_id = ?
			AND deleted = 0
			ORDER BY update_time
			LIMIT ?, ?`
	if err := mysql.DB.Select(&likes, sql, parentId, (current-1)*size, size); err != nil {
		return response.ERROR, likes
	}
	return response.SUCCESS, likes
}
