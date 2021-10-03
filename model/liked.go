package model

type Like struct {
	CreateTime    string `db:"create_time" json:"createTime"`
	UpdateTime    string `db:"update_time" json:"updateTime"`
	LikeId        string `db:"like_id" json:"likeId"`
	ParentId      string `db:"parent_id" json:"parentId"`
	LikeUser      string `db:"like_user" json:"likeUser"`
	LikePublisher string `db:"like_publisher" json:"likePublisher"`
	Deleted       int    `db:"deleted" json:"deleted"`
	PublisherInfo UserInfo
}
