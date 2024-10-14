package ds

import (
	"container/list"
	"errors"
	"sync"
)

// ref:https://dev.to/eyo000000/a-straightforward-guide-for-go-channel-3ba2
// 实现Go内建channel：queue存储数据、cond控制并发、closed状态

// Channel implement the Go channel
type Channel[T any] struct {
	queue    *list.List // store the data
	capacity int        // channel capacity
	cond     *sync.Cond // lock for synchronization
	closed   bool       // channel closed status
}

// NewChannel create a new channel
func NewChannel[T any](capacity int) *Channel[T] {
	return &Channel[T]{
		queue:    list.New(),
		capacity: capacity,
		cond:     sync.NewCond(&sync.Mutex{}),
		closed:   false,
	}
}

// Send data to channel
func (c *Channel[T]) Send(data T) error {
	// get lock for exclusive access
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	// return if channel is closed
	if c.closed {
		return errors.New("send on closed channel")
	}

	// wait until there is space in the channel
	for c.queue.Len() >= c.capacity && !c.closed {
		// if channel status is from normal to closed when this goroutine is sleeping,
		// we should check closed status again
		c.cond.Wait()
	}
	if c.closed {
		return errors.New("send on closed channel[x]")
	}
	// send data to channel
	c.queue.PushBack(data)

	// notify all waiting goroutines for receiving
	c.cond.Broadcast()

	// return success
	return nil
}

// Receive data from channel
// 当capacity为0时，send会阻塞，因为没有空闲位置可以存储数据；同时，receive会阻塞，因为没有数据可以读取。
// 具体是send阻塞在for c.queue.Len() >= c.capacity(=0)，receive阻塞在for c.queue.Len() == 0。
// send和receive都在等待对方接收数据和发送数据，因此会造成死锁!!!
//
// ref中用了一个很巧妙的方式解决了capacity=0时导致的死锁问题。capacity=0也就对应unbuffered channel(这非常有用)。
// 死锁不是 Send 和 Receive 由于锁竞争导致的，因为在等待时cond.Wait()会释放锁。死锁的根本原因是capacity=0，
// Send 和 Receive 主动阻塞而导致的死锁，因为我们定义了channel在capacity=0时，send和receive的行为是阻塞的。
// 为此，我们可以在receive时，先将capacity++，这样就能避免capacity=0时，send和receive阻塞。
// receive数据之后，再capacity--。
// 具体可以看 Receive 中特殊标记的line(########)。
func (c *Channel[T]) Receive() (T, bool) {
	// get lock for exclusive access
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	// return if channel is closed
	if c.closed {
		var zero T
		return zero, false
	}

	// ########## to solve deadlock when capacity = 0
	c.capacity++       // ++ to ensure there is space in the channel when capacity = 0
	c.cond.Broadcast() // notify sending goroutines which are waiting caused by capacity=0

	// wait until there is data in the channel
	for c.queue.Len() == 0 && !c.closed {
		// if channel status is from normal to closed when this goroutine is sleeping,
		// we should check closed status again
		c.cond.Wait()
	}
	if c.closed {
		var zero T
		return zero, false
	}

	// receive data from channel
	data := c.queue.Remove(c.queue.Front()).(T)

	// ########## to solve deadlock when capacity = 0
	c.capacity-- // recover capacity for correct action

	// notify all waiting goroutines for sending
	c.cond.Broadcast()

	// return data and success
	return data, true
}

// Close the channel
func (c *Channel[T]) Close() error {
	// get lock for exclusive access
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	if c.closed {
		return errors.New("close closed channel")
	}
	c.closed = true
	c.cond.Broadcast()
	return nil
}
