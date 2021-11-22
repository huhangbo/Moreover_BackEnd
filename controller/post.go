package controller

import (
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/post"
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
	code := post.PublishPost(tmpPost)
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
	code := post.UpdatePost(tmpPost)
	response.Response(c, code, nil)
}

func DeletePost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{
		PostId:    c.Param("postId"),
		Publisher: stuId.(string),
	}
	code := post.DeletePost(tmpPost)
	response.Response(c, code, nil)
}

func GetPostByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	switch c.Param("type") {
	case "page":
		code, posts, tmpPage := post.GetPostByPage(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"page":  tmpPage,
		})
	case "publisher":
		code, posts, tmpPage := post.GetPostByPublisher(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"page":  tmpPage,
		})
	case "top":
		code, posts, tmpPage := post.GetPostByTop(current, pageSize, stuId.(string))
		if code != response.SUCCESS {
			response.Response(c, code, nil)
			return
		}
		response.Response(c, code, gin.H{
			"posts": posts,
			"page":  tmpPage,
		})
	default:
		response.Response(c, response.ParamError, nil)
	}
}

func GetFollowPostByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, isEnd, posts := post.GetFollowPost(current, pageSize, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"posts": posts, "isEnd": isEnd})
}
