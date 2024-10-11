package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 实现带过期时间的单机锁，基于context实现
// ref:https://www.bilibili.com/video/BV1U14y167JP/?spm_id_from=333.999.0.0

type ExpireLock struct {
	coreLock    sync.Mutex
	processLock sync.Mutex

	stop  context.CancelFunc // 当持有锁的协程在锁还未过期之前主动释放锁时，通知异步解锁协程退出
	owner string             // 锁的持有者
}

func main() {
	locker := NewExpireLock()
	cnt := 0
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			time.Sleep(100 * time.Millisecond) // 等待所有协程都准备好了
			for j := 0; j < 10; j++ {
				locker.Lock(1)
				cnt++
				//locker.Unlock() // 即使没有主动释放锁，也不会导致死锁！！！
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(cnt)
}

func NewExpireLock() *ExpireLock {
	return &ExpireLock{
		coreLock:    sync.Mutex{},
		processLock: sync.Mutex{},
		stop:        nil,
		owner:       "",
	}
}

func (e *ExpireLock) Lock(expireMs int) {
	e.coreLock.Lock() // 如果锁还没有被释放，则会阻塞在这里
	// 保证后续的操作是无竞争的。即processLock的作用是保证加锁解锁的过程是串行的。
	e.processLock.Lock()
	defer e.processLock.Unlock()

	token := GetCurrentProcessAndGoroutineIdStr()
	e.owner = token
	if expireMs <= 0 {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	e.stop = cancel
	// 开一个异步协程，用于当锁过期后，释放锁
	go func() {
		select {
		case <-ctx.Done():
			// 走到此分支，说明是持有锁协程主动释放锁
			fmt.Printf("positive unlock, and async goroutine are announced to exit, token = %s\n", token)

		case <-time.After(time.Duration(expireMs) * time.Millisecond):
			// 走到此分支，说明是锁过期，需要异步协程释放锁
			fmt.Printf("lock expired, and released by the async goroutine, token = %s\n", token)
			_ = e.unlock(token)
		}
	}()
}

func (e *ExpireLock) Unlock() error {
	token := GetCurrentProcessAndGoroutineIdStr()
	return e.unlock(token)
}

// 会出现重复解锁的情况吗？
// 不会。每次解锁之前，会现持有processLock，保证解锁不会并发竞争。其次，解锁后会将owner置空，重复解锁的第二次或以后
// 的调用由于token != owner，表明已经结果锁，直接返回。
// 同理，持锁协程和异步解锁协程也不会重复解锁，因为解锁前是需要持有processLock。
func (e *ExpireLock) unlock(token string) error {
	e.processLock.Lock()
	defer e.processLock.Unlock()
	if token != e.owner {
		return errors.New("you are not lock owner")
	}
	if e.stop != nil {
		e.stop()
	}
	e.owner = ""
	e.coreLock.Unlock()
	return nil
}

func GetCurrentProcessAndGoroutineIdStr() string {
	pid := os.Getpid()
	return fmt.Sprintf("%d-%s", pid, GetCurrentGoroutineIdStr())
}

func GetCurrentGoroutineIdStr() string {
	buf := make([]byte, 128)
	n := runtime.Stack(buf, false)
	stackInfo := string(buf[:n])

	// 分割堆栈信息以找到goroutine ID
	parts := strings.Split(stackInfo, "\n")
	if len(parts) > 0 {
		firstLine := parts[0]
		if strings.Contains(firstLine, "goroutine") {
			// 提取 "goroutine 123" 中的数字部分
			idStr := strings.Fields(firstLine)[1]
			id, err := strconv.Atoi(idStr)
			if err == nil {
				return fmt.Sprintf("%d", id)
			}
		}
	}

	return ""
}
