package model

type UserInfo struct {
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

type UserInfoDetail struct {
	UserInfo
	Follower int
	Fan      int
	IsFollow bool
}
