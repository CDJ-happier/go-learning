package ds

import (
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	q := NewChannel[int](2)
	quit := make(chan struct{})
	go func() {
		for i := 0; i < 10; i++ {
			err := q.Send(i)
			if err != nil {
				fmt.Printf("send failed: %v\n", err)
			} else {
				fmt.Printf("send: %v\n", i)
			}
		}
		time.Sleep(1 * time.Second)
		err := q.Close()
		if err != nil {
			fmt.Printf("close failed: %v\n", err)
		}
		close(quit)
	}()

	for {
		v, ok := q.Receive()
		if !ok {
			break
		}
		fmt.Printf("receive: %v\n", v)
		time.Sleep(100 * time.Millisecond)
	}
	<-quit
}
