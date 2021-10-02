package cache

import "github.com/gin-gonic/gin"

func CacheHit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("cache", 0)
		c.Next()
		cache, _ := c.Get("cache")
		if cache.(string) == "1" {

		}
	}
}
