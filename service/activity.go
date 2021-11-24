package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

const (
	activityExpiration = time.Hour * 7 * 24
	sortActivityKey    = "activity:sort:"
)

func PublishActivity(activity dao.Activity) int {
	if err := conn.MySQL.Create(&activity).Error; err != nil {
		return response.FAIL
	}
	key := "activity:id:" + activity.ActivityId
	postJson, _ := json.Marshal(activity)
	pipe := conn.Redis.Pipeline()
	pipe.ZAdd(sortActivityKey, redis.Z{Member: activity.ActivityId, Score: float64(activity.CreatedAt.Unix())})
	pipe.ZAdd(sortActivityKey+activity.Category, redis.Z{Member: activity.ActivityId, Score: float64(activity.CreatedAt.Unix())})
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
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.ActivityId, "parent_id", stuId)
	return response.SUCCESS
}

func GetTotalByCategory(category string) (error, int64) {
	total, _ := conn.Redis.ZCard(sortActivityKey + category).Result()
	if total == 0 {
		var (
			tmpIds []struct {
				ActivityId string
				CreatedAt  time.Time
			}
			tmpZs []redis.Z
		)
		if category == "" {
			if err := conn.MySQL.Model(dao.Activity{}).Find(&tmpIds).Error; err != nil {
				return err, total
			}
		} else {
			if err := conn.MySQL.Model(&dao.Activity{}).Where("category = ?", category).Find(&tmpIds).Error; err != nil {
				return err, total
			}
		}
		for i := 0; i < len(tmpZs); i++ {
			tmpZs = append(tmpZs, redis.Z{Member: tmpIds[i], Score: float64(tmpIds[i].CreatedAt.Unix())})
		}
		if _, err := conn.Redis.ZAdd(sortActivityKey).Result(); err != nil {
			return err, total
		}
		return nil, int64(len(tmpZs))
	}
	return nil, total
}

func GetActivitiesByPublisher(current, size int, stuId string) (int, []dao.ActivityDetail, model.Page) {
	var (
		activities []dao.ActivityDetail
		ids        []string
		total      int64
	)
	if err := conn.MySQL.Model(&dao.Activity{}).Where("publisher = ?", stuId).Count(&total).Error; err != nil {
		return response.FAIL, activities, model.Page{}
	}
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if err := conn.MySQL.Model(&dao.Activity{}).Select("activity_id").Where("publisher = ?", stuId).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&ids).Error; err != nil {
		return response.FAIL, activities, tmpPage
	}
	for i := 0; i < len(ids); i++ {
		tmpActivityDetail := dao.ActivityDetail{Activity: dao.Activity{ActivityId: ids[i]}}
		if code := GetActivityDetailById(&tmpActivityDetail, stuId); code != response.SUCCESS {
			return code, activities, tmpPage
		}
		activities = append(activities, tmpActivityDetail)
	}
	return response.SUCCESS, activities, tmpPage
}

func GetActivitiesByCategory(current, size int, stuId, category string) (int, []dao.ActivityDetail, model.Page) {
	var (
		activities []dao.ActivityDetail
		tmpPage    model.Page
	)
	err, total := GetTotalByCategory(category)
	if err != nil {
		return response.FAIL, activities, tmpPage
	}
	tmpPage = model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: (int(total) / size) + 1}
	if (current-1)*size > int(total) {
		return response.PasswordError, activities, tmpPage
	}
	_, ids := util.GetIdsByPageFromRedis(current, size, "", "activity")
	for i := 0; i < len(ids); i++ {
		tmpActivityDetail := dao.ActivityDetail{Activity: dao.Activity{ActivityId: ids[i]}}
		if code := GetActivityDetailById(&tmpActivityDetail, stuId); code != response.SUCCESS {
			return code, activities, tmpPage
		}
		tmpActivityDetail.Detail = ""
		activities = append(activities, tmpActivityDetail)
	}
	return response.SUCCESS, activities, tmpPage
}
