package model

type Activity struct {
	PublishTime string `db:"create_time" json:"publishTime"`
	UpdateTime  string `db:"update_time" json:"updateTime"`
	ActivityId  string `db:"activity_id" json:"activityId"`
	Publisher   string `db:"publisher" json:"publisher"`
	Category    string `db:"category" json:"category"`
	Title       string `db:"title" json:"title"`
	Outline     string `db:"outline" json:"outline"`
	StartTime   string `db:"start_time" json:"startTime"`
	EndTime     string `db:"end_time" json:"endTime"`
	Contact     string `db:"contact" json:"contact"`
	Location    string `db:"location" json:"location"`
	Deleted     int    `db:"deleted" json:"deleted"`
	Star        int    `db:"star" json:"star"`
	Detail      string `db:"detail" json:"detail"`
}
