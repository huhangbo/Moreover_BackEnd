package util

import (
	"Moreover/conn"
)

func TopPost(postId, action string) error {
	var (
		key  = "post:sort:top"
		incr float64
	)
	switch action {
	case "liked":
		incr = 1
	case "dislike":
		incr = -1
	case "comment":
		incr = 3
	}
	return conn.Redis.ZIncrBy(key, incr, postId).Err()
}
