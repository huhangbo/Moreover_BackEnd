package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strconv"
)

func PublishPost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{PostId: uuid.New().String(), Publisher: stuId.(string)}
	if err := c.BindJSON(&tmpPost); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := service.PublishPost(tmpPost)
	response.Response(c, code, nil)
}

func UpdatePost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{
		PostId:    c.Param("postId"),
		Publisher: stuId.(string),
	}
	if err := c.BindJSON(&tmpPost); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := service.UpdatePost(tmpPost)
	response.Response(c, code, nil)
}

func DeletePost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{
		PostId:    c.Param("postId"),
		Publisher: stuId.(string),
	}
	code := service.DeletePost(tmpPost)
	response.Response(c, code, nil)
}

func GetPostByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	userId := c.Param("userId")
	switch c.Param("type") {
	case "page":
		code, posts, isEnd := service.GetPostByPage(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"isEnd": isEnd,
		})
	case "publisher":
		if userId == "" {
			userId = stuId.(string)
		}
		code, posts, isEnd := service.GetPostByPublisher(current, pageSize, stuId.(string), userId)
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"isEnd": isEnd,
		})
	case "top":
		code, posts, isEnd := service.GetPostByTop(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"isEnd": isEnd,
		})
	case "follow":
		code, posts, isEnd := service.GetFollowPost(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"isEnd": isEnd,
		})
	default:
		response.Response(c, response.ParamError, nil)
	}
}

func GetPostById(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	postId := c.Param("postId")
	PostDetail := dao.PostDetailFollow{
		Post: dao.Post{
			PostId: postId,
		},
	}
	code := service.GetPostDetailFollow(&PostDetail, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"content": PostDetail})
}
