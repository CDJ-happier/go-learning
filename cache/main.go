package main

import (
	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
)

// 源自: https://www.cyhone.com/articles/gin-cache/
func main() {
	ginCache()
}

// chenyahui/gin-cahce库的使用
func ginCache() {
	app := gin.New()

	memoryStore := persist.NewMemoryStore(1 * time.Minute)

	app.GET("/hello",
		cache.CacheByRequestURI(memoryStore, 2*time.Second),
		func(c *gin.Context) {
			time.Sleep(10 * time.Millisecond)
			c.String(200, "hello world")
		},
	)
	app.GET("/no_cache", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.String(200, "hello world")
	})

	if err := app.Run(":8080"); err != nil {
		panic(err)
	}
}
