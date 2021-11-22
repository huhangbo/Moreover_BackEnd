package message

import (
	"Moreover/conn"
	"Moreover/dao"
)

func PublishMessage(message dao.Message) error {
	return conn.MySQL.Create(&message).Error
}

func GetUnRead(action, stuId string) (error, int64) {
	var count int64
	if err := conn.MySQL.Model(dao.Message{}).Where("publisher = ? AND action = ? AND status = ?", stuId, action, 0).Count(&count).Error; err != nil {
		return err, count
	}
	return nil, count
}

func GetMessageByPage(current, size int, action, stuId string) (error, bool, []dao.Message) {
	var (
		messages []dao.Message
		isEnd    bool
	)
	if err := conn.MySQL.Model(dao.Message{}).Where("receiver = ? AND action = ?", stuId, action).Limit(size).Offset((current - 1) * size).Order("created_at DESC").Find(&messages).Error; err != nil {
		return err, isEnd, messages
	}
	if len(messages) < size {
		isEnd = true
	}
	return nil, isEnd, messages
}

func ReadAction(action, stuId string) error {
	return conn.MySQL.Model(dao.Message{}).Where("receiver = ? AND action = ? AND status = ?", stuId, action, 0).Update("status", 1).Error
}
