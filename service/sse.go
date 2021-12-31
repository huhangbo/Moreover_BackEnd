package service

import (
	"Moreover/dao"
	"sync"
)

const MsgQueLen = 3

var UserMap = userMap{Users: map[string]*UserData{}}

type userMap struct {
	Users   map[string]*UserData
	rwMutex sync.RWMutex
}

type UserData struct {
	StuId        string
	MessageQueue map[uint32]chan *dao.Message
	rwMutex      sync.RWMutex
}

func (t *UserData) InitMsgQue(serverId uint32) {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()
	t.MessageQueue[serverId] = make(chan *dao.Message, MsgQueLen)
}

func (t *UserData) SendMessage(message *dao.Message) {
	if message.Receiver == message.Publisher {
		return
	}
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()
	for serverId, msgQue := range t.MessageQueue {
		if len(msgQue) < MsgQueLen {
			msgQue <- message
		} else { //可能当前用户已下线channel被阻塞
			delete(t.MessageQueue, serverId)
		}
	}
}

func (t *UserData) GetMegQueue(severId uint32) <-chan *dao.Message {
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()
	return t.MessageQueue[severId]
}

func (t *userMap) PostMessage(message *dao.Message) {
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()
	if tmpUser, ok := t.Users[message.Receiver]; ok {
		go tmpUser.SendMessage(message)
	}
}

func (t *userMap) AddUser(stuId string, serverId uint32) *UserData {
	t.rwMutex.RLock()
	defer t.rwMutex.RUnlock()
	user, ok := t.Users[stuId]
	if ok {
		user.InitMsgQue(serverId)
	} else {
		user = &UserData{
			StuId:        stuId,
			MessageQueue: map[uint32]chan *dao.Message{serverId: make(chan *dao.Message, MsgQueLen)},
		}
		t.Users[stuId] = user
	}
	return user
}

func (t *userMap) RemoveUser(stuId string, serverId uint32) {
	t.rwMutex.Lock()
	defer t.rwMutex.Unlock()
	if len(t.Users[stuId].MessageQueue) == 1 { //只有一个连接
		delete(t.Users, stuId)
	} else { //多个设备连接
		delete(t.Users[stuId].MessageQueue, serverId)
	}
}
