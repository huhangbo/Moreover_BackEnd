package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var wg sync.WaitGroup

const (
	activityExpiration = time.Hour * 7 * 24
	sortActivityKey    = "activity:sort:"
)

func PublishActivity(activity dao.Activity) int {
	if err := conn.MySQL.Create(&activity).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	activity.PublishedAt = activity.CreatedAt.Unix()
	postJson, _ := json.Marshal(activity)
	pipe := conn.Redis.Pipeline()
	pipe.ZAdd(sortActivityKey, redis.Z{Member: activity.ActivityId, Score: float64(activity.PublishedAt)})
	pipe.ZAdd(sortActivityKey+activity.Category, redis.Z{Member: activity.ActivityId, Score: float64(activity.PublishedAt)})
	pipe.Set(key, string(postJson), activityExpiration)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func DeleteActivity(activity dao.Activity) int {
	tmpActivity := dao.Activity{ActivityId: activity.ActivityId}
	if code := GetActivityById(&tmpActivity); code != response.SUCCESS {
		return code
	}
	if err := conn.MySQL.Where("activity_id = ? AND publisher = ?", activity.ActivityId, activity.Publisher).Delete(&dao.Activity{}).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	pipe := conn.Redis.Pipeline()
	pipe.ZRem(sortActivityKey, activity.ActivityId)
	pipe.ZRem(sortActivityKey+tmpActivity.Category, activity.ActivityId)
	pipe.Del(key)
	if _, err := pipe.Exec(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func GetActivityById(activity *dao.Activity) int {
	key := "activity:id:" + activity.ActivityId
	activityString, err := conn.Redis.Get(key).Result()
	if err != nil {
		if err := conn.MySQL.Model(dao.Activity{}).Where("activity_id = ?", activity.ActivityId).First(activity).Error; err != nil {
			return response.FAIL
		}
		activity.PublishedAt = activity.CreatedAt.Unix()
		activityJson, _ := json.Marshal(activity)
		if _, err := conn.Redis.Set(key, string(activityJson), activityExpiration).Result(); err != nil {
			return response.FAIL
		}
	}
	_ = json.Unmarshal([]byte(activityString), activity)
	return response.SUCCESS
}

func UpdateActivity(activity dao.Activity) int {
	if err := conn.MySQL.Model(dao.Activity{}).Where("activity_id = ? AND publisher = ?", activity.ActivityId, activity.Publisher).Updates(activity).Error; err != nil {
		return response.FAIL
	}
	if err := conn.MySQL.Model(&dao.Activity{}).Where("activity_id = ?", activity.ActivityId).First(&activity).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	postJson, _ := json.Marshal(activity)
	if _, err := conn.Redis.Set(key, string(postJson), activityExpiration).Result(); err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func GetActivityDetailById(detail *dao.ActivityDetail, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: detail.ActivityId,
	}
	if code := GetActivityById(&tmpActivity); code != response.SUCCESS {
		return response.FAIL
	}
	detail.Activity = tmpActivity
	detail.PublisherInfo.StudentId = detail.Publisher
	GetUserInfoBasic(&(detail.PublisherInfo))
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.ActivityId, "parent", stuId)
	return response.SUCCESS
}

func GetActivityDetailFollow(detail *dao.ActivityDetailFollow, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: detail.ActivityId,
	}
	if code := GetActivityById(&tmpActivity); code != response.SUCCESS {
		return response.FAIL
	}
	detail.Activity = tmpActivity
	detail.PublisherInfo.StudentId = detail.Publisher
	GetUserInfoBasicFollow(&(detail.PublisherInfo), stuId)
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.ActivityId, "parent", stuId)
	return response.SUCCESS
}

func GetActivitiesByPublisher(current, size int, stuId, userId string) (int, []dao.ActivityDetail, bool) {
	var (
		activities []dao.ActivityDetail
		ids        []string
		isEnd      bool
	)
	if err := conn.MySQL.Model(&dao.Activity{}).Select("activity_id").Where("publisher = ?", userId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&ids).Error; err != nil {
		return response.FAIL, nil, isEnd
	}
	if len(ids) < size {
		isEnd = true
	}
	for i := 0; i < len(ids); i++ {
		tmpActivityDetail := dao.ActivityDetail{Activity: dao.Activity{ActivityId: ids[i]}}
		if code := GetActivityDetailById(&tmpActivityDetail, stuId); code != response.SUCCESS {
			return code, nil, isEnd
		}
		tmpActivityDetail.Detail = ""
		activities = append(activities, tmpActivityDetail)
	}
	return response.SUCCESS, activities, isEnd
}

func GetActivitiesByCategory(current, size int, stuId, category string) (int, []dao.ActivityDetail, bool) {
	var (
		activities []dao.ActivityDetail
		isEnd      bool
	)
	ids, _ := conn.Redis.ZRevRange(sortActivityKey+category, int64((current-1)*size), int64(current*size)-1).Result()
	if len(ids) == 0 {
		syncActivityToRedis(category)
		ids, _ = conn.Redis.ZRevRange(sortActivityKey+category, int64((current-1)*size), int64(current*size)-1).Result()
	}
	if len(ids) < size {
		isEnd = true
	}
	for i := 0; i < len(ids); i++ {
		tmpActivityDetail := dao.ActivityDetail{Activity: dao.Activity{ActivityId: ids[i]}}
		if code := GetActivityDetailById(&tmpActivityDetail, stuId); code != response.SUCCESS {
			return code, nil, isEnd
		}
		tmpActivityDetail.Detail = ""
		activities = append(activities, tmpActivityDetail)
	}
	return response.SUCCESS, activities, isEnd
}

func syncActivityToRedis(category string) {
	var (
		tmpIds []struct {
			ActivityId string
			CreatedAt  time.Time
		}
		tmpZs []redis.Z
	)
	if category == "" {
		if err := conn.MySQL.Model(dao.Activity{}).Find(&tmpIds).Error; err != nil {
			return
		}
	} else {
		if err := conn.MySQL.Model(&dao.Activity{}).Where("category = ?", category).Find(&tmpIds).Error; err != nil {
			return
		}
	}
	for i := 0; i < len(tmpZs); i++ {
		tmpZs = append(tmpZs, redis.Z{Member: tmpIds[i], Score: float64(tmpIds[i].CreatedAt.Unix())})
	}
	if _, err := conn.Redis.ZAdd(sortActivityKey).Result(); err != nil {
		return
	}
}
