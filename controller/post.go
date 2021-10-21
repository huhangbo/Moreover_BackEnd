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

func DeletePost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{
		PostId: c.Param("postId"),
	}
	code := post.DeletePost(tmpPost, stuId.(string))
	response.Response(c, code, nil)
}

func GetPostByPage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	code, posts, tmpPage := post.GetPostByPage(current, pageSize, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{
		"posts": posts,
		"page":  tmpPage,
	})
}

func UpdatePost(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	tmpPost := dao.Post{PostId: c.Param("postId")}
	if err := c.BindJSON(&tmpPost); err != nil {
		response.Response(c, response.ParamError, nil)
		return
	}
	code := post.UpdatePost(tmpPost, stuId.(string))
	response.Response(c, code, nil)
}
