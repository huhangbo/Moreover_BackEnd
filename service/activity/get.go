package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func GetActivityById(activity *dao.Activity) int {
	filter := bson.M{"_id": activity.ActivityId, "deleted": 0}
	if err := conn.MongoDB.Collection("activity").FindOne(context.TODO(), filter).Decode(&activity); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func GetActivityDetailById(detail *dao.ActivityDetail, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: detail.ActivityId,
	}
	code := GetActivityById(&tmpActivity)
	detail.Activity = tmpActivity
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.ActivityId, "parent_id", stuId)

	return code
}

func GetTotal(filter bson.M) (int, int64) {
	count, err := conn.MongoDB.Collection("activity").CountDocuments(context.TODO(), filter)
	if err != nil {
		return response.ParamError, count
	}
	return response.SUCCESS, count
}
