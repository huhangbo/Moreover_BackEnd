package liked

import (
	"Moreover/conn"
	"Moreover/dao"
	"Moreover/pkg/response"
	"Moreover/service/util"
	"time"
)

func PublishLike(liked dao.Liked) int {
	if err := conn.MySQL.Create(&liked).Error; err != nil {
		return response.FAIL
	}
	if !util.PublishSortRedis(liked.Publisher, float64(liked.CreatedAt.Unix()), "liked:sort:"+liked.Parent, "liked:tmp:"+liked.Parent) {
		return response.FAIL
	}
	return response.SUCCESS
}

func SyncLike2MySQL() {
	for {
		ticker := time.NewTicker(time.Hour * 12)
		<-ticker.C
	}

}

func StartTime() {
	go func() {
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		t := time.NewTimer(next.Sub(now))
		<-t.C
	}()
}
