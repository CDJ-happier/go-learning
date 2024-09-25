package ds

import (
	"fmt"
	"testing"
)

func TestLRUCache(t *testing.T) {
	capacity := 10
	lruCache := NewLRUCache(capacity)
	for i := 0; i < capacity; i++ {
		lruCache.Add(fmt.Sprintf("key.%d", i), i)
	}
	t.Log(lruCache.Keys())
	for i := 0; i < capacity; i += 2 {
		lruCache.Add(fmt.Sprintf("key.%d", i), i)
	}
	t.Log(lruCache.Keys())
}
