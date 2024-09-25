package ds

import "sync"

// 双向链表节点
type lruNode struct {
	key  string
	val  interface{}
	next *lruNode
	prev *lruNode
}

type LRUCache struct {
	capacity int
	size     int
	head     *lruNode
	keys     map[string]*lruNode
	mu       sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	head := &lruNode{
		key:  "head",
		val:  nil,
		next: nil,
		prev: nil,
	}
	head.next = head
	head.prev = head
	return &LRUCache{
		capacity: capacity,
		size:     0,
		head:     head,
		keys:     make(map[string]*lruNode),
	}
}

// 将node从链表中脱离
func (l *LRUCache) removeNode(node *lruNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// 将node插入到after之后
func (l *LRUCache) addNodeAfter(node *lruNode, after *lruNode) {
	node.next = after.next
	after.next = node
	node.prev = after
	node.next.prev = node
}

func (l *LRUCache) Add(key string, val interface{}) {
	// key是否存在
	l.mu.Lock()
	defer l.mu.Unlock()
	node, ok := l.keys[key]
	if ok {
		node.val = val
		l.removeNode(node)
	} else {
		node = &lruNode{
			key:  key,
			val:  val,
			next: nil,
			prev: nil,
		}
		l.size++
		if l.size > l.capacity {
			// 删除链表尾部节点
			tail := l.head.prev
			l.removeNode(tail)
			delete(l.keys, tail.key)
			l.size--
		}
		l.keys[key] = node
	}
	// 放到链表头部
	l.addNodeAfter(node, l.head)
}

func (l *LRUCache) Get(key string) interface{} {
	l.mu.Lock()
	defer l.mu.Unlock()
	node, ok := l.keys[key]
	if !ok {
		return nil
	}
	// 放到链表头部
	l.removeNode(node)
	l.addNodeAfter(node, l.head)
	return node.val
}

func (l *LRUCache) Delete(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	node, ok := l.keys[key]
	if !ok {
		return
	}
	l.removeNode(node)
	delete(l.keys, key)
	l.size--
}

func (l *LRUCache) IsEmpty() bool {
	return l.size == 0
}

func (l *LRUCache) Size() int {
	return l.size
}

func (l *LRUCache) Capacity() int {
	return l.capacity
}

func (l *LRUCache) Keys() []string {
	keys := make([]string, 0, l.size)
	node := l.head.next
	for node != l.head {
		keys = append(keys, node.key)
		node = node.next
	}
	return keys
}
