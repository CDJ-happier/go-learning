[toc]

# Go Repo

[使用singlefilght解决缓存击穿](./singleflight_test.go)

[令牌桶限流](./ratelimiter/main.go)

golang中通过`time.Now()`返回当前日期时间结构体`Time`，可以通过其方法`Year(),Month()`等获取具体日期、时间，且`Time`还提供一些加减方法。

`time.Now().Unix()`返回自Unix Epoch起，以毫秒为单位的时间戳。可以通过`time.Unix(timestamp)`将时间戳转为`Time`。

golang中的`time.Duration()`是一个int64的类型，本质表示两个时间点之间以**纳秒**为单位的间隔，可以通过其方法`Seconds()`之类转为秒(返回float64类型)。

[布隆过滤器](./bloomfilter/main.go)

一般设置插入的个数n、允许的错误率p。通过这两个参数可以计算出bloom filter的bit大小m，以及hash个数k。公式如下：
$$
m = \frac{-n{\cdot}\ln^{p}}{(\ln^{2})^2} \\
k = \frac{m}{n} {\cdot} \ln^{2}
$$
当n=300, p=0.1时，m=1438， k=4。

一个简单hash的设计：

1. 先将key映射为一个值base；
2. 针对第i个hash，我们加上一个偏移，如base += i * 101；
3. 最后index = base % m；

[基于context实现带过期时间的单机锁](./expirelock/main.go)

context用于超时解锁。

[手写Go channel](./ds/channel.go)

要求：并发安全、queue-like action、send、receive、unbuffered&buffered。[ref](https://dev.to/eyo000000/a-straightforward-guide-for-go-channel-3ba2)

除了基础实现外，ref中通过一个很巧妙的方式解决capacity=0(unbuffered channel)时导致的deadlock。由于unbuffered，send没有空间写，因此睡眠、receive因为没有数据读，因此睡眠，导致了死锁。一个解决方案就是在receive时，将capacity++，然后cond.Notify通知所有因capacity=0导致的send goroutine，然后读取后再将capacity自减。

```go
// Channel implement the Go channel
type Channel[T any] struct {
	queue    *list.List // store the data
	capacity int        // channel capacity
    // wait when no space to send for sending goroutine
    // wait when no data to receive for receive goroutine
	cond     *sync.Cond // lock, and producer&consumer architecture
	closed   bool       // channel closed status
}
```



## Data Structure

[LRU](./ds/lru.go)：使用一个双向链表（带虚拟头节点）+哈希表实现。哈希表获得O(1)时间复杂度的查询，双向链表用于维护节点的使用时间（逻辑上头节点的next是最新使用的）

[LFU](./ds/lfu.go)：使用一个哈希表实现，key=freq，val=该freq下的节点所组成的链表。每次使用一个节点元素时，从map[freq]链表中移除该元素，并加入到map[freq+1]的链表中。

[优先队列](./ds/priority_queue.go)：基于堆实现，利用数组存储元素，主要有`Push,Pop,Top,Size`几个方法。插入元素时，在尾部插入，并与父节点比较是否满足堆结构要求，否则交换（Up）。删除元素时，是删除堆顶元素，然后把堆底元素放到堆顶，并与子节点比较是否满足对结构要求，否则交换（Down）。

使用数组保存堆元素的一些性质：堆大小为size。非叶子节点i满足0≤i≤size/2-1，节点i的left=2i+1, right=2i+2，parent=i/2-1。

[跳表](./ds/skiplist.go)：一种概率数据结构，空间换时间。

## Algorithm

[排序算法](./algorithm/sort.go)：快速排序、堆排序（堆化、下沉）