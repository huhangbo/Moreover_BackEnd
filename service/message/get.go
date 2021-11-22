package message

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLatest(action, stuId, current string) (int, dao.Message, int64) {
	var message dao.Message
	filter := bson.M{"action": action, "receiver": stuId, "read": 0}
	err := conn.MongoDB.Collection("message").FindOne(context.TODO(), filter).Decode(&message)
	count, _ := conn.MongoDB.Collection("message").CountDocuments(context.TODO(), filter)
	if err != nil {
		return response.FAIL, message, count
	}
	return response.SUCCESS, message, count
}
