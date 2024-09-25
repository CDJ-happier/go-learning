package main

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

// RateLimiter 限流器接口定义
type RateLimiter[T comparable] interface {
	Take(ctx context.Context, item T) bool
	When(ctx context.Context, item T) time.Duration
	Forget(ctx context.Context, item T)
	Retries(ctx context.Context, item T) int
}

type TokenRateLimiter[T comparable] struct {
	burst      int
	rate       float64
	tokens     int
	lastUpdate int64 // Unix timestamp, in milliseconds since 1970-01-01
	mu         sync.Mutex
}

func NewTokenRateLimiter[T comparable](burst int, rate float64) RateLimiter[T] {
	return &TokenRateLimiter[T]{
		burst:      burst,
		rate:       rate,
		tokens:     burst,
		lastUpdate: time.Now().Unix(),
	}
}

func (rl *TokenRateLimiter[T]) AllowN(n int) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	duration := now.Sub(time.Unix(rl.lastUpdate, 0))
	durationS := duration.Seconds()
	durationS = max(durationS, 0)
	curCap := rl.tokens + int(durationS*rl.rate)
	curCap = min(curCap, rl.burst)
	if curCap >= n {
		rl.tokens = curCap - n
		rl.lastUpdate = now.Unix()
		return true
	}
	rl.tokens = curCap
	rl.lastUpdate = now.Unix()
	return false
}

func (rl *TokenRateLimiter[T]) WaitN(n int) time.Duration {
	if rl.AllowN(n) {
		return time.Duration(0)
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()
	waitTime := time.Duration(int64(float64(n-rl.tokens)/rl.rate)) * time.Second
	return waitTime
}

func (rl *TokenRateLimiter[T]) Take(_ context.Context, _ T) bool {
	return rl.AllowN(1)
}

func (rl *TokenRateLimiter[T]) When(_ context.Context, _ T) time.Duration {
	return rl.WaitN(1)
}

func (rl *TokenRateLimiter[T]) Forget(_ context.Context, _ T) {
}

func (rl *TokenRateLimiter[T]) Retries(_ context.Context, _ T) int {
	return 0
}

// ItemExponentialFailureRateLimiter 指数退避限流
// wait = base * 2^(failures) <= maxDelay
type ItemExponentialFailureRateLimiter[T comparable] struct {
	items     map[T]int
	baseDelay time.Duration
	maxDelay  time.Duration
	mu        sync.Mutex
}

func NewItemExponentialFailureRateLimiter[T comparable](baseDelay, maxDelay time.Duration) RateLimiter[T] {
	return &ItemExponentialFailureRateLimiter[T]{
		items:     make(map[T]int),
		baseDelay: baseDelay,
		maxDelay:  maxDelay,
	}
}

func (irl *ItemExponentialFailureRateLimiter[T]) Take(_ context.Context, item T) bool {
	irl.mu.Lock()
	defer irl.mu.Unlock()
	_, ok := irl.items[item]
	return !ok
}

func (irl *ItemExponentialFailureRateLimiter[T]) When(_ context.Context, item T) time.Duration {
	irl.mu.Lock()
	defer irl.mu.Unlock()
	failures, ok := irl.items[item]
	if !ok {
		irl.items[item] = 1
		return irl.baseDelay
	}
	delay := irl.baseDelay * time.Duration(math.Pow(2, float64(failures)))
	if delay > irl.maxDelay {
		delay = irl.maxDelay
	}
	irl.items[item]++
	return delay
}

func (irl *ItemExponentialFailureRateLimiter[T]) Forget(_ context.Context, item T) {
	irl.mu.Lock()
	defer irl.mu.Unlock()
	delete(irl.items, item)
}
func (irl *ItemExponentialFailureRateLimiter[T]) Retries(_ context.Context, item T) int {
	return irl.items[item]
}

func main() {
	limiter := NewTokenRateLimiter[string](10, 2)
	tokenLimiter := limiter.(*TokenRateLimiter[string])
	for i := 0; i < 10; i++ {
		wait := tokenLimiter.WaitN(4)
		if wait > 0 {
			fmt.Printf("wait %v\n", wait)
			time.Sleep(wait)
		} else {
			fmt.Printf("allowN\n")
		}
	}
}

func parallelRequest() {
	limiter := NewTokenRateLimiter[string](10, 2)
	tokenLimiter := limiter.(*TokenRateLimiter[string])
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(name string) {
			for i := 0; i < 10; i++ {
				if tokenLimiter.AllowN(2) {
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
