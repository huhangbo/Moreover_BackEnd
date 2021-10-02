package model

type User struct {
	StudentID string `db:"student_id"`
	UserName  string `db:"username"`
	Password  string `db:"password"`
	Deleted   bool   `db:"deleted"`
}
