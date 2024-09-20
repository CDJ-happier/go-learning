package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type RateLimiter struct {
	burst      int
	rate       int
	capacity   int
	lastUpdate int64 // Unix timestamp, in milliseconds since 1970-01-01
	mu         sync.Mutex
}

func NewRateLimiter(burst int, rate int) *RateLimiter {
	return &RateLimiter{
		burst:      burst,
		rate:       rate,
		capacity:   burst,
		lastUpdate: time.Now().Unix(),
	}
}

func (rl *RateLimiter) AllowN(n int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	duration := now.Sub(time.Unix(rl.lastUpdate, 0))
	durationS := duration.Seconds()
	durationS = max(durationS, 0)
	curCap := rl.capacity + int(durationS*float64(rl.rate))
	curCap = min(curCap, rl.burst)
	if curCap >= n {
		rl.capacity = curCap - n
		rl.lastUpdate = now.Unix()
		return true
	}
	rl.capacity = curCap
	rl.lastUpdate = now.Unix()
	return false
}

func (rl *RateLimiter) WaitN(n int) time.Duration {
	if rl.AllowN(n) {
		return time.Duration(0)
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()
	waitTime := time.Duration((n-rl.capacity)/rl.rate) * time.Second
	return waitTime
}

func main() {
	limiter := NewRateLimiter(10, 2)
	for i := 0; i < 10; i++ {
		wait := limiter.WaitN(4)
		if wait > 0 {
			fmt.Printf("wait %v\n", wait)
			time.Sleep(wait)
		} else {
			fmt.Printf("allowN\n")
		}
	}
}

func parallelRequest() {
	limiter := NewRateLimiter(10, 1)
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(name string) {
			for i := 0; i < 10; i++ {
				if limiter.AllowN(2) {
					fmt.Printf("%s allowN\n", name)
				} else {
					fmt.Printf("%s Not allowN\n", name)
				}
				time.Sleep(1 * time.Second)
			}
			wg.Done()
		}("user" + strconv.FormatInt(int64(i), 10))
	}
	wg.Wait()
}
