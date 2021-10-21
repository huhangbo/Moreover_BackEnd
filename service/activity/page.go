package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/liked"
	"Moreover/service/user"
	"github.com/go-redis/redis"
	"time"
)

func GetActivitiesByPade(current, size int, category, stuId string) (int, []dao.ActivityBasic, model.Page) {
	var activities []dao.ActivityBasic
	var tmpActivities []dao.Activity
	var activityIds []string
	code, total := GetTotal(category)
	var tmpPage = model.Page{
		Current:   current,
		PageSize:  size,
		Total:     total,
		TotalPage: (total / size) + 1,
	}
	if code != response.SUCCESS || (current-1)*size > total {
		return response.ParamError, activities, tmpPage
	}
	code, activityIds = getActivityIdsByPageFromRedis(current, size, category)
	code, tmpActivities = getActivityByIds(activityIds)
	if code != response.SUCCESS || (len(activityIds) == 0 && code == response.SUCCESS) {
		if err := conn.MySQL.Limit(size).Offset((current - 1) * size).Find(&tmpActivities).Error; err != nil {
			return response.FAIL, activities, tmpPage
		}
	}
	for i := 0; i < len(tmpActivities); i++ {
		tmpActivityBasic := getActivityBasicByActivity(tmpActivities[i], stuId)
		activities = append(activities, tmpActivityBasic)
	}
	return response.SUCCESS, activities, tmpPage
}

func getActivityBasicByActivity(activity dao.Activity, stuId string) dao.ActivityBasic {
	tmpBasic := dao.ActivityBasic{
		CreatedAt:  activity.CreatedAt,
		UpdatedAt:  activity.UpdatedAt,
		ActivityId: activity.ActivityId,
		Publisher:  activity.Publisher,
		Category:   activity.Category,
		Title:      activity.Title,
		Outline:    activity.Outline,
		StartTime:  activity.StartTime,
		EndTime:    activity.EndTime,
		Location:   activity.Location,
		PublisherInfo: dao.UserInfoBasic{
			StudentId: activity.Publisher,
		},
	}
	user.GetUserInfoBasic(&(tmpBasic.PublisherInfo))
	_, tmpBasic.Star, tmpBasic.IsStar = liked.GetTotalAndLiked(tmpBasic.ActivityId, stuId)
	return tmpBasic
}

func getActivityIdsByPageFromRedis(current, size int, category string) (int, []string) {
	sortKey := "activity:sort:"
	if category != "" {
		sortKey = "activity:sort:" + category
	}
	activitiesId, err := conn.Redis.ZRevRange(sortKey, int64((current-1)*size), int64(current*size)).Result()
	if err != nil {
		return response.ERROR, nil
	}
	return response.SUCCESS, activitiesId
}

func GetActivitiesByPublisher(current, size int, stuId string) (int, []dao.ActivityBasic, model.Page) {
	var tmpActivitiesBasic []dao.ActivityBasic
	key := "activity:publisher:" + stuId
	total, _ := conn.Redis.ZCard(key).Result()
	tmpPage := model.Page{
		Current:   current,
		PageSize:  size,
		Total:     int(total),
		TotalPage: (int(total) / size) + 1,
	}
	if total == 0 {
		if !publisherActivityRedis(stuId) {
			return response.SUCCESS, tmpActivitiesBasic, tmpPage
		}
		total, _ = conn.Redis.ZCard(key).Result()
	}
	tmpPage.Total = int(total)
	tmpPage.TotalPage = (int(total) / size) + 1
	if (current-1)*size > int(total) {
		return response.ParamError, tmpActivitiesBasic, tmpPage
	}
	activityIds, _ := conn.Redis.ZRevRange(key, int64((current-1)*size), int64(current*size)).Result()
	code, tmpActivities := getActivityByIds(activityIds)
	if code != response.SUCCESS {
		return code, tmpActivitiesBasic, tmpPage
	}
	for i := 0; i < len(tmpActivities); i++ {
		tmpActivityBasic := getActivityBasicByActivity(tmpActivities[i], stuId)
		tmpActivitiesBasic = append(tmpActivitiesBasic, tmpActivityBasic)
	}
	return response.SUCCESS, tmpActivitiesBasic, tmpPage
}

func publisherActivityRedis(stuId string) bool {
	var tmpActivities []dao.Activity
	conn.MySQL.Select("activity_id, created_at").Where("publisher = ?", stuId).Find(&tmpActivities)
	key := "activity:publisher:" + stuId
	pipe := conn.Redis.Pipeline()
	for _, item := range tmpActivities {
		pipe.ZAdd(key, redis.Z{Member: item, Score: float64(item.CreatedAt.Unix())})
	}
	pipe.Expire(key, time.Minute*5)

	if _, err := pipe.Exec(); err != nil || len(tmpActivities) == 0 {
		return false
	}
	return true
}
