package model

type UserInfo struct {
	StudentID   string `db:"student_id" json:"student_id"`
	Nickname    string `db:"nickname" json:"nickname"`
	Sex         string `db:"sex" json:"sex"`
	Description string `db:"description" json:"description"`
	Avatar      string `db:"avatar" json:"avatar"`
	Tag         string `db:"tag" json:"tag"`
}
