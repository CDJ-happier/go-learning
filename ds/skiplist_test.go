package ds

import (
	"fmt"
	"testing"
)

func TestSkipList_PutGet(t *testing.T) {
	sl := NewSkipList()
	for i := 0; i < 100; i++ {
		sl.Put(i, fmt.Sprintf("hi,%d", i))
	}
	for i := 0; i < 100; i++ {
		if sl.Get(i) != fmt.Sprintf("hi,%d", i) {
			t.FailNow()
		}
		fmt.Printf("%d:%s\n", i, sl.Get(i))
	}
}

func TestSkipList_Delete(t *testing.T) {
	sl := NewSkipList()
	for i := 0; i < 100; i++ {
		sl.Put(i, fmt.Sprintf("hi,%d", i))
	}
	for i := 20; i < 80; i++ {
		sl.Delete(i)
	}
	cur := sl.head.next[0]
	for cur != nil {
		fmt.Printf("%d:%s\n", cur.key, cur.val)
		cur = cur.next[0]
	}
}

func TestSkipList_Range(t *testing.T) {
	sl := NewSkipList()
	for i := 0; i < 100; i++ {
		sl.Put(i, fmt.Sprintf("hi,%d", i))
	}
	for _, v := range sl.Range(20, 80) {
		fmt.Printf("%s\n", v)
	}
}

func TestSkipList_CeilingFloor(t *testing.T) {
	sl := NewSkipList()
	for i := 0; i < 100; i++ {
		sl.Put(i, fmt.Sprintf("hi,%d", i))
	}
	fmt.Printf("%s\n", sl.Ceiling(-1))
	fmt.Printf("%s\n", sl.Floor(200))
}

func BenchmarkSkipList_PutGet(b *testing.B) {
	sl := NewSkipList()
	for i := 0; i < b.N; i++ {
		sl.Put(i, fmt.Sprintf("hi,%d", i))
		sl.Get(i)
	}
}
