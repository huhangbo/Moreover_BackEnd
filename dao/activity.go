package dao

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type Activity struct {
	ID          uint                  `gorm:"autoIncrement primaryKey" json:"-"`
	PublishedAt int64                 `gorm:"-" json:"publishedAt"`
	CreatedAt   time.Time             `json:"-"`
	UpdatedAt   time.Time             `json:"-"`
	DeletedAt   soft_delete.DeletedAt `json:"-"`
	ActivityId  string                `json:"activityId"`
	Publisher   string                `json:"publisher"`
	Category    string                `json:"category"`
	Title       string                `json:"title"`
	Outline     string                `json:"outline"`
	StartTime   int                   `json:"startTime"`
	EndTime     int                   `json:"endTime"`
	Location    string                `json:"location"`
	Detail      string                `json:"detail"`
	Contact     string                `json:"contact"`
}

type ActivityDetail struct {
	Activity
	Star          int           `json:"star"`
	IsStar        bool          `json:"isStar"`
	PublisherInfo UserInfoBasic `json:"publisherInfo"`
}

type ActivityDetailFollow struct {
	Activity
	Star          int                 `json:"star"`
	IsStar        bool                `json:"isStar"`
	PublisherInfo UserInfoBasicFollow `json:"publisherInfo"`
}
