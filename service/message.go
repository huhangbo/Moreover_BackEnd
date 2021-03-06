package service

import (
	"Moreover/conn"
	"Moreover/dao"
	"errors"
)

func PublishMessage(message dao.Message) error {
	if message.Publisher == message.Receiver {
		return errors.New("repeat")
	}
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
	for i := 0; i < len(messages); i++ {
		messages[i].PublishedAt = messages[i].CreatedAt.Unix()
	}
	return nil, isEnd, messages
}

func ReadAction(action, stuId string) error {
	return conn.MySQL.Model(dao.Message{}).Where("receiver = ? AND action = ? AND status = ?", stuId, action, 0).Update("status", 1).Error
}
