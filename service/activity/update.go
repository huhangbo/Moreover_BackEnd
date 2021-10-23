package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateActivity(activity dao.Activity) int {
	tmpActivity := dao.Activity{
		ActivityId: activity.ActivityId,
	}
	GetActivityById(&tmpActivity)
	if tmpActivity.Publisher != activity.Publisher {
		return response.AuthError
	}
	tmp, _ := bson.Marshal(activity)
	var tmpBson bson.M
	if err := bson.Unmarshal(tmp, &tmpBson); err != nil {
		return response.FAIL
	}
	update := bson.M{"$set": tmpBson}
	if _, err := conn.MongoDB.Collection("activity").UpdateOne(context.TODO(), bson.M{"_id": activity.ActivityId, "deleted": 0}, update); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
