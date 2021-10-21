package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Activity struct {
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
	DeletedAt  soft_delete.DeletedAt `json:"deletedAt"`
	ActivityId string                `gorm:"primaryKey" json:"activityId"`
	Publisher  string                `json:"publisher"`
	Category   string                `json:"category"`
	Title      string                `json:"title"`
	Outline    string                `json:"outline"`
	StartTime  string                `json:"start_time"`
	EndTime    string                `json:"end_time"`
	Location   string                `json:"location"`
	Detail     string                `json:"detail"`
	Contact    string                `json:"contact"`
}

type ActivityBasic struct {
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ActivityId    string    `json:"activityId"`
	Publisher     string    `json:"publisher"`
	Category      string    `json:"category"`
	Title         string    `json:"title"`
	Outline       string    `json:"outline"`
	StartTime     string    `json:"startTime"`
	EndTime       string    `json:"endTime"`
	Location      string    `json:"location"`
	Star          int       `json:"star"`
	IsStar        bool      `json:"isSar"`
	PublisherInfo UserInfoBasic
}

type ActivityDetail struct {
	Activity
	Star          int  `json:"star"`
	IsStar        bool `json:"isStar"`
	PublisherInfo UserInfoBasic
}
