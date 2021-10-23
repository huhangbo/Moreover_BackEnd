package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"context"
)

func PublishActivity(activity dao.Activity) int {
	if _, err := conn.MongoDB.Collection("activity").InsertOne(context.TODO(), activity); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
