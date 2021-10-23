package dao

import (
	"time"
)

type Activity struct {
	CreatedAt  time.Time `bson:"created_at,omitempty" json:"createdAt"`
	Deleted    int       `bson:"deleted" json:"-"`
	ActivityId string    `bson:"_id" json:"activityId"`
	Publisher  string    `json:"publisher"`
	Category   string    `json:"category"`
	Title      string    `json:"title"`
	Outline    string    `json:"outline"`
	StartTime  string    `bson:"start_time" json:"startTime"`
	EndTime    string    `bson:"end_time" json:"endTime"`
	Location   string    `json:"location"`
	Detail     string    `json:"detail"`
	Contact    string    `json:"contact"`
}

type ActivityBasic struct {
	CreatedAt     time.Time `bson:"created_at,omitempty" json:"createdAt"`
	Deleted       int       `bson:"deleted" json:"-"`
	ActivityId    string    `bson:"_id" json:"activityId"`
	Publisher     string    `json:"publisher"`
	Category      string    `json:"category"`
	Title         string    `json:"title"`
	Outline       string    `json:"outline"`
	StartTime     string    `bson:"start_time" json:"startTime"`
	EndTime       string    `bson:"end_time" json:"endTime"`
	Location      string    `json:"location"`
	Contact       string    `json:"contact"`
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
