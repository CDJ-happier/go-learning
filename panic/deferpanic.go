package main

import (
	"fmt"
	"runtime"
	"time"
)

/*
// doCall handles the single call for a key.
func (g *Group) doCall(c *call, key string, fn func() (interface{}, error)) {
	normalReturn := false
	recovered := false

	// use double-defer to distinguish panic from runtime.Goexit,
	// more details see https://golang.org/cl/134395
	defer func() {
		// the given function invoked runtime.Goexit
		if !normalReturn && !recovered {
			c.err = errGoexit
		}

		g.mu.Lock()
		defer g.mu.Unlock()
		c.wg.Done()
		if g.m[key] == c {
			delete(g.m, key)
		}

		if e, ok := c.err.(*panicError); ok {
			// In order to prevent the waiting channels from being blocked forever,
			// needs to ensure that this panic cannot be recovered.
			if len(c.chans) > 0 {
				go panic(e)
				select {} // Keep this goroutine around so that it will appear in the crash dump.
			} else {
				panic(e)
			}
		} else if c.err == errGoexit {
			// Already in the process of goexit, no need to call again
		} else {
			// Normal return
			for _, ch := range c.chans {
				ch <- Result{c.val, c.err, c.dups > 0}
			}
		}
	}()

	// 1. 执行匿名函数
	func() {
		defer func() {
			if !normalReturn {
				// Ideally, we would wait to take a stack trace until we've determined
				// whether this is a panic or a runtime.Goexit.
				//
				// Unfortunately, the only way we can distinguish the two is to see
				// whether the recover stopped the goroutine from terminating, and by
				// the time we know that, the part of the stack trace relevant to the
				// panic has been discarded.
				if r := recover(); r != nil {
					c.err = newPanicError(r)
				}
				// 如果recover返回的不是nil，则说明fn中发生了panic。此时，通过recover进行处理并记录panic。因为使用了recover，
				//	因此，对于doCall来说，该匿名函数是正常的，会执行最后三行以标记发生了panic。
				// 如果返回的是nil，则说明fn中发生了Goexit。此时会向上传递，即对于doCall来说，该匿名函数发生了Goexit，
				// 会执行doCall的defer，最后三行不会执行，因此normalReturn和recovered都是false。因此，在doCall的defer中会对
				// c.err标记为errGoexit。
				// 我觉得把对recovered的赋值写在上面newPanicError后面更好理解。
			}
		}()
		// 2. 如果fn中发生了panic或者Goexit的调用（且panic没有recover），则会在defer中进行处理。
		c.val, c.err = fn()
		normalReturn = true
	}()

	if !normalReturn {
		recovered = true
	}
}
*/

// 该案例源自理解singleflight中doCall的double defer区分panic和runtime.Goexit。原理见匿名函数的注释（recover的返回值是不是nil）
func main() {
	fmt.Println("demo1 #####################")
	demo1()
	fmt.Println()
	fmt.Println("demo2 #####################")
	demo2()
	fmt.Println()
	fmt.Println("demo3 #####################")
	// runtime.Goexit()与panic的区别：？
	// runtime.Goexit()的应用场景？
	demoGoexit()
	time.Sleep(100 * time.Millisecond)
	// Goexit()的标准库注释
	// Goexit terminates the goroutine that calls it. No other goroutine is affected.
	// Goexit runs all deferred calls before terminating the goroutine. Because Goexit
	// is not a panic, any recover calls in those deferred functions will return nil.
	//
	// Calling Goexit from the main goroutine terminates that goroutine
	// without func main returning. Since func main has not returned,
	// the program continues execution of other goroutines.
	// If all other goroutines exit, the program crashes.
}

func demo1() {
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from panic outer:", r)
			}
		}()
		// 因为在 invokePanicNoRecover() 中 panic 了，并且没有使用 recover()处理，那么panic()会继续向上抛出。
		// 因此，对于当前匿名函数来说，invokePanicNoRecover()是一个panic，会停止执行后续的语句并执行当前匿名函数的defer 。
		invokePanicNoRecover()
		fmt.Println("statement after invokePanicNoRecover() won't execute")
	}()
}

func demo2() {
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from panic outer:", r)
			}
		}()
		// 因为在 invokePanicWithRecover() 中 panic 了，并且使用了 recover()处理，那么panic()会停止向上抛出。
		// 因此，对于当前匿名函数来说，invokePanicWithRecover()是一个正常的函数，会继续正常执行。
		invokePanicWithRecover()
		fmt.Println("statement after invokePanicWithRecover() will execute")
	}()
}

func demoGoexit() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from panic:", r)
			} else {
				fmt.Println("no need to recover, Goexit")
			}
		}()
		invokeGoexit()
		fmt.Println("statement after invokeGoexit() won't execute")
	}()
}

func invokePanicNoRecover() {
	defer fmt.Println("output after panic() by defer")
	panic("my panic")
	//fmt.Println("execution after panic()") // 无法到达
}

func invokePanicWithRecover() {
	defer func() {
		fmt.Println("output after panic() by defer")
		if r := recover(); r != nil {
			fmt.Println("recovered from panic inner:", r)
		}
	}()
	panic("my panic")
}

func invokeGoexit() {
	defer fmt.Println("output after goexit() by defer")
	runtime.Goexit()
	fmt.Println("execution after goexit()") // unreachable
}
