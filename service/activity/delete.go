package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func DeleteActivity(activity dao.Activity, stuId string) int {
	GetActivityById(&activity)
	if activity.Publisher != stuId {
		return response.AuthError
	}
	update := bson.M{"$set": bson.M{"deleted": 1}}
	if _, err := conn.MongoDB.Collection("activity").UpdateOne(context.TODO(), bson.M{"_id": activity.ActivityId}, update); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
