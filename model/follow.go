package model

type Follow struct {
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
	Deleted    int    `db:"deleted" json:"deleted"`
	Follower   string `db:"follower" json:"follower"`
	Fan        string `db:"fan" json:"fan"`
}
