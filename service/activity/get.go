package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/user"
	"Moreover/service/util"
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
)

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

func GetActivityDetailById(detail *dao.ActivityDetail, stuId string) int {
	tmpActivity := dao.Activity{
		ActivityId: detail.ActivityId,
	}
	if code := GetActivityById(&tmpActivity); code != response.SUCCESS {
		return response.FAIL
	}
	detail.Activity = tmpActivity
	detail.PublisherInfo.StudentId = detail.Publisher
	user.GetUserInfoBasic(&(detail.PublisherInfo))
	_, detail.Star, detail.IsStar = util.GetTotalAndIs("liked", detail.ActivityId, "parent_id", stuId)
	return response.SUCCESS
}

func GetTotalByCategory(category string) (error, int64) {
	total, _ := conn.Redis.ZCard(sortKey + category).Result()
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
		if _, err := conn.Redis.ZAdd(sortKey).Result(); err != nil {
			return err, total
		}
		return nil, int64(len(tmpZs))
	}
	return nil, total
}
