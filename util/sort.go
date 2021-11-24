package util

import (
	"Moreover/conn"
	"github.com/go-redis/redis"
	"time"
)

type Sort struct {
	UpdatedAt time.Time `json:"updatedAt"`
	id        string
}

func PublishSortRedis(id string, score float64, keys ...string) bool {
	pipe := conn.Redis.Pipeline()
	value := redis.Z{
		Member: id,
		Score:  score,
	}
	for _, item := range keys {
		pipe.ZAdd(item, value)
		pipe.Expire(item, time.Hour*24*7)
	}
	if _, err := pipe.Exec(); err != nil {
		return false
	}
	return true
}

func DeleteSortRedis(id, key string, keys ...string) bool {
	pipe := conn.Redis.Pipeline()
	for _, item := range keys {
		pipe.ZRem(item, id)
	}
	pipe.Del(key)
	if _, err := pipe.Exec(); err != nil {
		return false
	}
	return true
}
