package model

type UserInfo struct {
	//StudentID   string `db:"student_id" json:"student_id"`
	//	//Nickname    string `db:"nickname" json:"nickname"`
	//	//Sex         string `db:"sex" json:"sex"`
	//	//Avatar      string `db:"avatar" json:"avatar"`
	//	//Tag         string `db:"tag" json:"tag"`
	UserBasicInfo
	Description string `db:"description" json:"description"`
}

type UserBasicInfo struct {
	StudentID string `db:"student_id" json:"student_id"`
	Nickname  string `db:"nickname" json:"nickname"`
	Avatar    string `db:"avatar" json:"avatar"`
	Sex       string `db:"sex" json:"sex"`
	Tag       string `db:"tag" json:"tag"`
}
