package model

type Comment struct {
	PublishTime string `db:"create_time" json:"publishTime"`
	UpdateTime  string `db:"update_time" json:"updateTime"`
	CommentId   string `db:"comment_id" json:"commentId"`
	ParentID    string `db:"parent_id" json:"parentId"`
	Publisher   string `db:"publisher" json:"publisher"`
	Deleted     int    `db:"deleted" json:"deleted"`
	Star        int    `db:"star" json:"star"`
	Message     string `db:"message" json:"message"`
}
