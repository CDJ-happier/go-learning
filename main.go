package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	limiter := rate.NewLimiter(rate.Every(time.Second), 10)
	fmt.Println("after limiter")
	for i := 0; i < 100; i++ {
		fmt.Println(limiter.Allow())
		time.Sleep(20 * time.Millisecond)
	}
}
