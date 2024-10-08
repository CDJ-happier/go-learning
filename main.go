package main

import (
	"fmt"
	"go-repo/ds"
)

func main() {
	capacity := 10
	lruCache := ds.NewLRUCache(capacity)
	for i := 0; i < capacity; i++ {
		lruCache.Add(fmt.Sprintf("key.%d", i), i)
	}
	fmt.Println(lruCache.Keys())
	for i := 0; i < capacity; i += 2 {
		lruCache.Add(fmt.Sprintf("key.%d", i), i+2)
	}
	fmt.Println(lruCache.Keys())
	arr := make([]string, 0)
	for i := 0; i < capacity; i++ {
		arr = append(arr, fmt.Sprintf("key.%d", i))
	}
	fmt.Println(lruCache.Get(arr[0]))
}
