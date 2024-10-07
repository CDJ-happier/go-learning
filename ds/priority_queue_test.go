package ds

import (
	"fmt"
	"testing"
)

func TestNewPriorityQueue(t *testing.T) {
	pq := NewPriorityQueue[int](10, func(a, b int) bool {
		return a < b // 小根堆
	})
	data := []int{5, 3, 4, 1, 2}
	for _, v := range data {
		pq.Push(v)
		fmt.Printf("top element is: %v, size=%v\n", pq.Top(), pq.Size())
	}
}

func TestPriorityQueue_IsEmpty(t *testing.T) {
	pq := NewPriorityQueue[int](10, func(a, b int) bool {
		return a < b
	})
	fmt.Println(pq.IsEmpty())
	pq.Push(1)
	fmt.Println(pq.IsEmpty())
}
