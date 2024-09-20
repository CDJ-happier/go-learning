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

