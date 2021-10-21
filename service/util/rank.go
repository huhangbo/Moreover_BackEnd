package util

import "time"

func GetTopScore(star int, createdAt time.Time) float64 {
	hour := time.Now().Hour() - createdAt.Hour()
	return float64((star+1)/(hour+2) ^ 2)
}
