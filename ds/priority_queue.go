package ds

type LessFunc[T comparable] func(T, T) bool

type PriorityQueue[T comparable] struct {
	data []T
	less LessFunc[T]
	size int
}

// 基于顺序存储的堆
// 堆顶元素为data[0]，堆底元素为data[size-1]，堆大小为size
// 非叶子节点的索引为i，其中0 <= i <= size/2-1。
// 节点i的左子节点为2i+1，右子节点为2i+2，父节点为(i-1)/2

func NewPriorityQueue[T comparable](capacity int, less LessFunc[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		data: make([]T, capacity),
		less: less,
		size: 0,
	}
}

func (pq *PriorityQueue[T]) Push(x T) {
	// 将元素x添加到末尾，并与父节点比较是否满足堆结构要求，如果不满足则交换
	pq.data[pq.size] = x
	pq.size++
	pq.siftUp(pq.size - 1)
}

func (pq *PriorityQueue[T]) siftUp(i int) {
	for {
		parent := (i - 1) / 2
		if parent < 0 || !pq.less(pq.data[i], pq.data[parent]) {
			// 如果已经Up到父节点或者当前节点i满足堆结构要求，退出循环
			break
		}
		pq.data[i], pq.data[parent] = pq.data[parent], pq.data[i]
		i = parent
	}
}

func (pq *PriorityQueue[T]) Pop() T {
	// 返回堆顶元素。并将末尾元素移动到堆顶，并重新调整堆结构
	x := pq.data[0]
	pq.data[0] = pq.data[pq.size-1]
	pq.size--
	pq.siftDown(0)
	return x
}

func (pq *PriorityQueue[T]) siftDown(i int) {
	for {
		left := 2*i + 1
		right := 2*i + 2
		smaller := i
		if left < pq.size && pq.less(pq.data[left], pq.data[smaller]) {
			smaller = left
		}
		if right < pq.size && pq.less(pq.data[right], pq.data[smaller]) {
			smaller = right
		}
		if smaller != i {
			pq.data[i], pq.data[smaller] = pq.data[smaller], pq.data[i]
			i = smaller
		} else {
			break
		}
	}
}

func (pq *PriorityQueue[T]) Top() T {
	return pq.data[0]
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return pq.size == 0
}

func (pq *PriorityQueue[T]) Size() int {
	return pq.size
}
