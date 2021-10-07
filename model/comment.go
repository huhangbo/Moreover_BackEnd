package model

type Comment struct {
	CreateTime string `db:"create_time" json:"createTime"`
	UpdateTime string `db:"update_time" json:"updateTime"`
	Deleted    int    `db:"deleted" json:"deleted"`
	CommentId  string `db:"comment_id" json:"commentId"`
	ParentID   string `db:"parent_id" json:"parentId"`
	Publisher  string `db:"publisher" json:"publisher"`
	Replier    string `db:"replier"   json:"replier"`
	Star       int    `json:"star"`
	Message    string `db:"message" json:"message"`
}

type CommentDetail struct {
	Comment
	PublisherInfo UserInfo
}

type ChildComment struct {
	Comments []CommentDetail
	Total    int `json:"total"`
}

type ParentComment struct {
	CommentDetail
	Children ChildComment
}
