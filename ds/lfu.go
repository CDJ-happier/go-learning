package ds

import "container/list"

// LFUCache
// https://leetcode.cn/problems/lfu-cache/solutions/2457716/tu-jie-yi-zhang-tu-miao-dong-lfupythonja-f56h/

type entry struct {
	key, value, freq int
}

type LFUCache struct {
	capacity   int
	minFreq    int
	keyToNode  map[int]*list.Element
	freqToList map[int]*list.List
}

func NewLFUCache(capacity int) LFUCache {
	return LFUCache{
		capacity:   capacity,
		keyToNode:  map[int]*list.Element{},
		freqToList: map[int]*list.List{},
	}
}

func (l *LFUCache) pushFront(e *entry) {
	if _, ok := l.freqToList[e.freq]; !ok {
		l.freqToList[e.freq] = list.New() // 双向链表
	}
	l.keyToNode[e.key] = l.freqToList[e.freq].PushFront(e)
}

func (l *LFUCache) getEntry(key int) *entry {
	node := l.keyToNode[key]
	if node == nil { // 没有这本书
		return nil
	}
	e := node.Value.(*entry)
	lst := l.freqToList[e.freq]
	lst.Remove(node)    // 把这本书抽出来
	if lst.Len() == 0 { // 抽出来后，这摞书是空的
		delete(l.freqToList, e.freq) // 移除空链表
		if l.minFreq == e.freq {     // 这摞书是最左边的
			l.minFreq++
		}
	}
	e.freq++       // 看书次数 +1
	l.pushFront(e) // 放在右边这摞书的最上面
	return e
}

func (l *LFUCache) Get(key int) int {
	if e := l.getEntry(key); e != nil {
		return e.value
	}
	return -1
}

func (l *LFUCache) Put(key, value int) {
	if e := l.getEntry(key); e != nil {
		e.value = value // 更新 value
		return
	}
	if len(l.keyToNode) == l.capacity {
		lst := l.freqToList[l.minFreq]
		delete(l.keyToNode, lst.Remove(lst.Back()).(*entry).key)
		if lst.Len() == 0 {
			delete(l.freqToList, l.minFreq)
		}
	}
	l.pushFront(&entry{key, value, 1})
	l.minFreq = 1
}
