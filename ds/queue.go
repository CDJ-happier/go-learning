package ds

type Queue[T comparable] struct {
	items []T
	size  int
}

func NewQueue[T comparable]() *Queue[T] {
	return &Queue[T]{
		items: make([]T, 0),
		size:  0,
	}
}

func (q *Queue[T]) Push(x T) {
	q.items = append(q.items, x)
	q.size++
}

func (q *Queue[T]) Pop() T {
	if q.size == 0 {
		panic("queue is empty")
	}
	x := q.items[0]
	q.items = q.items[1:]
	q.size--
	return x
}

func (q *Queue[T]) Front() T {
	if q.size == 0 {
		panic("queue is empty")
	}
	return q.items[0]
}

func (q *Queue[T]) Size() int {
	return q.size
}

func (q *Queue[T]) IsEmpty() bool {
	return q.size == 0
}
