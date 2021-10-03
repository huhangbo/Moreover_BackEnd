package model

type Comment struct {
	PublishTime   string `db:"create_time" json:"publishTime"`
	UpdateTime    string `db:"update_time" json:"updateTime"`
	CommentId     string `db:"comment_id" json:"commentId"`
	ParentID      string `db:"parent_id" json:"parentId"`
	Publisher     string `db:"publisher" json:"publisher"`
	Replier       string `db:"replier"   json:"replier"`
	Deleted       int    `db:"deleted" json:"deleted"`
	Star          int    `json:"star"`
	Message       string `db:"message" json:"message"`
	PublisherInfo UserInfo
	Children      CommentList
}

type CommentList struct {
	Comments []Comment
	Page     `json:"page"`
}
