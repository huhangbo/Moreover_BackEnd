package model

type Base struct {
	PublishTime string `db:"create_time" json:"publishTime"`
	UpdateTime  string `db:"update_time" json:"updateTime"`
	Deleted     int    `db:"deleted" json:"deleted"`
}
