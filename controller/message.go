package controller

import (
	"Moreover/pkg/response"
	"Moreover/service/message"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"sync/atomic"
)

var messageKind = []string{"comment", "liked", "follow"}

var ServerId uint32

func HandleSSE(c *gin.Context) {
	stu, _ := c.Get("stuId")
	stuId := stu.(string)
	current := c.Param("current")
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	serverId := atomic.AddUint32(&ServerId, 1)
	tmpUser := message.UserMap.AddUser(stuId, serverId)
	messageChan := tmpUser.GetMegQueue(serverId)
	for i := 0; i < len(messageKind); i++ {
		item := messageKind[i]
		go func() {
			_, tmpMessage, tmpCount := message.GetLatest("comment", stuId, current)
			c.SSEvent(item, gin.H{"message": tmpMessage, "count": tmpCount})
		}()
	}
	disConn := c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-messageChan:
			c.SSEvent(msg.Kind, gin.H{"message": msg})
			return true
		}
	})
	if disConn {
		message.UserMap.RemoveUser(stuId, serverId)
	}
}

func GetMessages(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	action := c.Param("action")
	code, messages, tmpPage := message.GetMessageByPage(int64(current), int64(pageSize), action, stuId.(string))
	if code != response.SUCCESS {
		response.Response(c, code, nil)
		return
	}
	response.Response(c, code, gin.H{"message": messages, "page": tmpPage})
}

func ReadMessage(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	action := c.Param("action")
	code := message.ReadAction(action, stuId.(string))
	response.Response(c, code, nil)
}
