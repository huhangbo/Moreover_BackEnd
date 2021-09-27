package model

type User struct {
	StudentID string `db:"student_id"`
	UserName string `db:"username"`
	Password string	`db:"password"`
	Telephone string `db:"telephone"`
	Deleted bool `db:"deleted"`
}