package activity

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
	"Moreover/service/util"
)

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
