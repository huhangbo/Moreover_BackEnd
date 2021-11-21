package message

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
)

func PublishMessage(message dao.Message) int {
	if err := conn.MySQL.Create(&message).Error; err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
