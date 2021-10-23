package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetActivitiesByPublisher(current, size int64, stuId string) (int, []dao.Activity, model.Page) {
	filter := bson.M{"deleted": 0, "publisher": stuId}
	var activities []dao.Activity
	code, total := GetTotal(filter)
	skip := (current - 1) * size
	option := &options.FindOptions{Limit: &size, Skip: &skip}
	var tmpPage = model.Page{
		Current:   int(current),
		PageSize:  int(size),
		Total:     int(total),
		TotalPage: int((total / size) + 1),
	}
	if code != response.SUCCESS || (current-1)*size > total {
		return code, activities, tmpPage
	}
	results, _ := conn.MongoDB.Collection("activity").Find(context.TODO(), filter, option)
	if err := results.All(context.TODO(), &activities); err != nil {
		return response.FAIL, nil, tmpPage
	}
	return response.SUCCESS, activities, tmpPage
}

func GetActivitiesByPade(current, size int64, stuId, category string) (int, []dao.ActivityBasic, model.Page) {
	var activities []dao.ActivityBasic
	filter := bson.M{"deleted": 0, "category": category}
	if category == "" {
		filter = bson.M{"deleted": 0}
	}
	code, total := GetTotal(filter)
	var tmpPage = model.Page{
		Current:   int(current),
		PageSize:  int(size),
		Total:     int(total),
		TotalPage: int((total / size) + 1),
	}
	if code != response.SUCCESS || (current-1)*size > total {
		return code, activities, tmpPage
	}
	skip := (current - 1) * size
	option := &options.FindOptions{Limit: &size, Skip: &skip}
	results, _ := conn.MongoDB.Collection("activity").Find(context.TODO(), filter, option)
	if err := results.All(context.TODO(), &activities); err != nil {
		return response.FAIL, nil, tmpPage
	}
	for i := 0; i < len(activities); i++ {
		activities[i].PublisherInfo.StudentId = activities[i].Publisher
		user.GetUserInfoBasic(&(activities[i].PublisherInfo))
		_, activities[i].Star, activities[i].IsStar = liked.GetTotalAndLiked(activities[i].ActivityId, stuId)
	}
	return response.SUCCESS, activities, tmpPage
}
