package message

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"context"
)

func PublishMessage(message dao.Message) int {

	if _, err := conn.MongoDB.Collection("message").InsertOne(context.TODO(), message); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
