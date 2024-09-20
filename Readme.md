# Go惯用法

[使用singlefilght解决缓存击穿](./singleflight_test.go)

[令牌桶限流](./ratelimiter/main.go)

golang中通过`time.Now()`返回当前日期时间结构体`Time`，可以通过其方法`Year(),Month()`等获取具体日期、时间，且`Time`还提供一些加减方法。

`time.Now().Unix()`返回自Unix Epoch起，以毫秒为单位的时间戳。可以通过`time.Unix(timestamp)`将时间戳转为`Time`。

golang中的`time.Duration()`是一个int64的类型，本质表示两个时间点之间以**纳秒**为单位的间隔，可以通过其方法`Seconds()`之类转为秒(返回float64类型)。

