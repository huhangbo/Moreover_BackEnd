package controller

import (
	"Moreover/pkg/response"
	"Moreover/service"
	"github.com/gin-gonic/gin"
	"io"
	"strconv"
	"sync/atomic"
)

var (
	messageKind = []string{"comment", "liked", "follow"}
	ServerId    uint32
)

func HandleSSE(c *gin.Context) {
	stu, _ := c.Get("stuId")
	stuId := stu.(string)
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	serverId := atomic.AddUint32(&ServerId, 1)
	tmpUser := service.UserMap.AddUser(stuId, serverId)
	messageChan := tmpUser.GetMegQueue(serverId)
	for i := 0; i < len(messageKind); i++ {
		item := messageKind[i]
		go func() {
			_, count := service.GetUnRead(item, stuId)
			if count != 0 {
				c.SSEvent(item, gin.H{"count": count})
			}
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
		service.UserMap.RemoveUser(stuId, serverId)
	}
}

func GetMessages(c *gin.Context) {
	stuId, _ := c.Get("stuId")
	current, _ := strconv.Atoi(c.Param("current"))
	pageSize, _ := strconv.Atoi(c.Param("pageSize"))
	action := c.Param("action")
	err, isEnd, messages := service.GetMessageByPage(current, pageSize, action, stuId.(string))
	if err != nil {
		response.Response(c, response.FAIL, nil)
		return
	}
	if err := service.ReadAction(action, stuId.(string)); err != nil {
		response.Response(c, response.FAIL, nil)
		return
	}
	response.Response(c, response.SUCCESS, gin.H{"messages": messages, "isEnd": isEnd})
}
