package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"testing"
	"time"
)

// Response 模拟数据库返回的数据结构
type Response struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Id   int    `json:"id"`
}

var count int
var group singleflight.Group

// GetUserInfo 模拟从数据库获取用户信息
// 一般是先从缓存获取，缓存没有再从数据库获取. 如果发生缓存击穿，并且此时有多个请求同时发起，那么就会导致数据库压力过大
// 解决方案是使用分布式锁，加锁线程从数据库获取数据并写入缓存，其他请求等待锁释放，再从缓存获取数据
func GetUserInfo(name string) (Response, error) {
	count++ // 模拟多个请求并发访问时的耗时操作(对数据库的压力递增)
	time.Sleep(time.Duration(count) * time.Millisecond)
	return Response{
		Name: name,
		Age:  18,
		Id:   1,
	}, nil
}

// GetUserInfoBundle 使用singleflight.Group解决缓存击穿。原理同分布式锁，但singleflight.Group是单机实现。
func GetUserInfoBundle(sf *singleflight.Group, name string) (Response, error) {
	resp, err, _ := sf.Do(name, func() (interface{}, error) {
		return GetUserInfo(name)
	})
	return resp.(Response), err
}

func helper(getType string) {
	wg := sync.WaitGroup{}
	start := time.Now()
	numRequests := 1000
	key := "dijkstraCai"
	count = 0
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			//var resp Response
			var err error
			if getType == "singleflight" {
				_, err = GetUserInfoBundle(&group, key)
			}
			if getType == "concurrency" {
				_, err = GetUserInfo(key)
			}
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Printf("resp: %v\n", resp)
		}(key)
	}
	wg.Wait()
	fmt.Printf("time elapse for all requests is: %v\n", time.Since(start))
}

// TestSingleflightGet 测试singleflight.Group
func TestSingleflightGet(t *testing.T) {
	helper("singleflight")
}

// TestConcurrencyGet 测试并发访问
func TestConcurrencyGet(t *testing.T) {
	helper("concurrency")
}
