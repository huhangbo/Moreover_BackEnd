package model

import "time"

type Model struct {
	ID          int64 `db:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deleted     bool
	DeletedTime time.Time
}