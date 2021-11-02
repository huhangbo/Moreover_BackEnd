package message

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMessageByPage(current, size int64, action, stuId string) (int, []dao.Message, model.Page) {
	var messages []dao.Message
	filter := bson.M{"action": action, "receiver": stuId}
	_, total := GetTotal(action, stuId)
	tmpPage := model.Page{Current: int(current), PageSize: int(size), Total: int(total), TotalPage: int((total / size) + 1)}
	if (current-1)*size > total {
		return response.ParamError, messages, tmpPage
	}
	skip := (current - 1) * size
	option := &options.FindOptions{Limit: &size, Skip: &skip, Sort: bson.M{"created_at": -1}}
	result, _ := conn.MongoDB.Collection("message").Find(context.TODO(), filter, option)
	if err := result.All(context.TODO(), &messages); err != nil {
		return response.FAIL, messages, tmpPage
	}
	return response.SUCCESS, messages, tmpPage
}

func GetTotal(action, stuId string) (int, int64) {
	filter := bson.M{"action": action, "receiver": stuId}
	count, err := conn.MongoDB.Collection("message").CountDocuments(context.TODO(), filter)
	if err != nil {
		return response.FAIL, count
	}
	return response.SUCCESS, count
}

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

func ReadAction(action, stuId string) int {
	filter := bson.M{"action": action, "receiver": stuId}
	update := bson.M{"$set": bson.M{"read": 1}}
	_, err := conn.MongoDB.Collection("message").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
