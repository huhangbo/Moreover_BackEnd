package dao

import "time"

type Message struct {
	CreateAt  time.Time `bson:"create_at"`
	Read      int
	Parent    string
	Kind      string
	Action    string
	Receiver  string
	Publisher string
	Detail    string
}
