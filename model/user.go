package model

type User struct {
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
	Deleted    int    `db:"deleted" json:"deleted"`
	StudentID  string `db:"student_id"`
	UserName   string `json:"username"`
	Password   string `json:"password"`
}
