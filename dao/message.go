package dao

import "time"

type Message struct {
	PublishedAt int64     `gorm:"-" json:"publishedAt"`
	ID          uint      `gorm:"autoIncrement primaryKey" json:"-"`
	CreatedAt   time.Time `json:"-"`
	Receiver    string    `json:"-"`
	Publisher   string    `json:"publisher"`
	Status      int       `json:"status"`
	Parent      string    `json:"parent"`
	Kind        string    `json:"kind"`
	Action      string    `json:"action"`
	Detail      string    `json:"detail"`
}
