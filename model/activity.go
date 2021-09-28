package model

type Activity struct {
	PublishTime string `db:"publish_time"`
	UpdateTime  string `db:"update_time"`
	ActivityId  string `db:"activity_id"`
	Publisher   string `db:"publisher"`
	Category    string `db:"category" json:"category"`
	Title       string `db:"title" json:"title"`
	Outline     string `db:"outline" json:"outline"`
	StartTime   string `db:"start_time" json:"startTime"`
	EndTime     string `db:"end_time" json:"endTime"`
	Contact     string `db:"contact" json:"contact"`
	Location    string `db:"location" json:"location"`
	Detail      string `db:"detail" json:"detail"`
}
