package message

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/model"
	"Moreover/pkg/response"
)

func PublishMessage(message dao.Message) int {
	if err := conn.MySQL.Create(&message).Error; err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}

func GetUnRead(action, stuId string) (int, []dao.Message) {
	var messages []dao.Message
	if err := conn.MySQL.Model(dao.Message{}).Where("publisher = ? AND action = ? AND status = ?", stuId, action, 0).Find(&messages).Error; err != nil {
		return response.FAIL, messages
	}
	return response.SUCCESS, messages
}

func GetTotal(action, stuId string) (int, int64) {
	var total int64
	if err := conn.MySQL.Model(dao.Message{}).Where("publisher = ? AND action = ?", stuId, action).Count(&total).Error; err != nil {
		return response.FAIL, total
	}
	return response.SUCCESS, total
}

func GetMessageByPage(current, size int, action, stuId string) (int, []dao.Message, model.Page) {
	var messages []dao.Message
	_, total := GetTotal(action, stuId)
	tmpPage := model.Page{Current: current, PageSize: size, Total: int(total), TotalPage: int(total)/size + 1}
	if (current-1)*size > int(total) {
		return response.ParamError, messages, tmpPage
	}
	if err := conn.MySQL.Model(dao.Message{}).Where("publisher = ? AND action = ?", stuId, action).Find(&messages).Error; err != nil {
		return response.FAIL, messages, tmpPage
	}
	return response.SUCCESS, messages, tmpPage
}

func ReadAction(action, stuId string) int {
	if err := conn.MySQL.Model(dao.Message{}).Where("publisher = ? AND action = ? AND status = ?", stuId, action, 0).Update("status = ?", 1).Error; err != nil {
		return response.FAIL
	}
	return response.SUCCESS
}
