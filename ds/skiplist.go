package ds

import (
	"math/rand"
)

type skipNode struct {
	key  int
	val  interface{}
	next []*skipNode
}

type SkipList struct {
	head *skipNode
}

// NewSkipList ...
func NewSkipList() *SkipList {
	return &SkipList{
		head: &skipNode{
			key:  -1,
			next: make([]*skipNode, 0),
		},
	}
}

// randomLevel 以一定概率生成level
func randomLevel() int {
	level := 0
	// 每次都是1/2的概率
	for rand.Intn(2) == 0 {
		level++
	}
	return level
}

func (s *SkipList) search(key int) *skipNode {
	cur := s.head
	for level := len(cur.next) - 1; level >= 0; level-- {
		for cur.next[level] != nil && cur.next[level].key < key {
			cur = cur.next[level]
		}
		if cur.next[level] != nil && cur.next[level].key == key {
			return cur.next[level]
		}
	}
	return nil
}

func (s *SkipList) Get(key int) interface{} {
	if n := s.search(key); n != nil {
		return n.val
	}
	return nil
}

func (s *SkipList) Put(key int, val interface{}) {
	if n := s.search(key); n != nil {
		n.val = val
	}
	curLevel := randomLevel()
	newNode := &skipNode{
		key:  key,
		val:  val,
		next: make([]*skipNode, curLevel+1),
	}
	// 如果新节点的level大于当前链表的level，则需要增加链表的level
	for len(s.head.next)-1 < curLevel {
		s.head.next = append(s.head.next, nil)
	}
	// 在每层插入节点
	cur := s.head
	for level := curLevel; level >= 0; level-- {
		for cur.next[level] != nil && cur.next[level].key < key {
			cur = cur.next[level]
		}
		// 在当前level插入节点
		newNode.next[level] = cur.next[level]
		cur.next[level] = newNode
	}
}

func (s *SkipList) Delete(key int) {
	if n := s.search(key); n == nil {
		return
	}
	cur := s.head
	for level := len(cur.next) - 1; level >= 0; level-- {
		for cur.next[level] != nil && cur.next[level].key < key {
			cur = cur.next[level]
		}
		if cur.next[level] != nil && cur.next[level].key == key {
			cur.next[level] = cur.next[level].next[level]
		}
	}
	// 删除有些节点后会导致level减少，需要清理
	diff := 0
	for level := len(s.head.next) - 1; level >= 0 && s.head.next[level] == nil; level-- {
		diff++
	}
	s.head.next = s.head.next[:len(s.head.next)-diff]
}

// ceiling 返回大于等于key的最小值
func (s *SkipList) ceiling(key int) *skipNode {
	cur := s.head
	for level := len(cur.next) - 1; level >= 0; level-- {
		for cur.next[level] != nil && cur.next[level].key < key {
			cur = cur.next[level]
		}
		if cur.next[level] != nil && cur.next[level].key == key {
			return cur.next[level]
		}
	}
	return cur.next[0]
}

// floor 返回小于等于key的最大值
func (s *SkipList) floor(key int) *skipNode {
	cur := s.head
	for level := len(cur.next) - 1; level >= 0; level-- {
		for cur.next[level] != nil && cur.next[level].key < key {
			cur = cur.next[level]
		}
		if cur.next[level] != nil && cur.next[level].key == key {
			return cur.next[level]
		}
	}
	return cur
}

func (s *SkipList) Ceiling(key int) interface{} {
	if n := s.ceiling(key); n != nil {
		return n.val
	}
	return nil
}

func (s *SkipList) Floor(key int) interface{} {
	if n := s.floor(key); n != nil {
		return n.val
	}
	return nil
}

func (s *SkipList) Range(start, end int) []interface{} {
	n := s.ceiling(start)
	if n == nil {
		return nil
	}
	var res []interface{}
	for ; n != nil && n.key <= end; n = n.next[0] {
		res = append(res, n.val)
	}
	return res
}
