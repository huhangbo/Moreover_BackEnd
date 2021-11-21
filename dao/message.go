package dao

import "time"

type Message struct {
	ID        uint      `gorm:"autoIncrement primaryKey" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"-"`
	Read      int       `json:"read"`
	Parent    string    `json:"parent"`
	Kind      string    `json:"kind"`
	Action    string    `json:"action"`
	Receiver  string    `json:"receiver"`
	Publisher string    `json:"publisher"`
	Detail    string    `json:"detail"`
}
