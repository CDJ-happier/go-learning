package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	ErrCacheMiss = errors.New("cache miss")
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

type Cacher interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, val interface{}, expire time.Duration) error
}

type CacheProxy struct {
	localCache  Cacher
	remoteCache Cacher
}

func NewCacheProxy(localCache Cacher, remoteCache Cacher) *CacheProxy {
	return &CacheProxy{localCache: localCache, remoteCache: remoteCache}
}

func (p *CacheProxy) Get(ctx context.Context, key string) (interface{}, error) {
	if p.localCache == nil {
		return p.remoteCache.Get(ctx, key)
	}
	val, err := p.localCache.Get(ctx, key)
	if err == nil {
		return val, nil
	}
	if !errors.Is(ErrCacheMiss, err) {
		return nil, err
	}
	val, err = p.remoteCache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	err = p.localCache.Set(ctx, key, val, 0)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (p *CacheProxy) Set(ctx context.Context, key string, val interface{}, expire time.Duration) error {
	var err error
	err = p.remoteCache.Set(ctx, key, val, expire)
	if err != nil {
		return fmt.Errorf("remote cache set failed: %w", err)
	}
	if p.localCache != nil {
		return p.localCache.Set(ctx, key, val, expire)
	}
	return nil
}
