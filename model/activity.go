package model

type Activity struct {
	CreateTime    string `db:"create_time" json:"createTime"`
	UpdateTime    string `db:"update_time" json:"updateTime"`
	ActivityId    string `db:"activity_id" json:"activityId"`
	Publisher     string `db:"publisher" json:"publisher"`
	Category      string `db:"category" json:"category"`
	Title         string `db:"title" json:"title"`
	Outline       string `db:"outline" json:"outline"`
	StartTime     string `db:"start_time" json:"startTime"`
	EndTime       string `db:"end_time" json:"endTime"`
	Contact       string `db:"contact" json:"contact"`
	Location      string `db:"location" json:"location"`
	Deleted       int    `db:"deleted" json:"deleted"`
	Star          int    `json:"star"`
	Detail        string `db:"detail" json:"detail"`
	PublisherInfo UserInfo
}
