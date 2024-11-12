package osutil

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestAfterStopFunc(t *testing.T) {
	var lastIndex int
	AfterStopFunc(10*time.Second, func(c <-chan struct{}) {
		for i := 5; i >= 0; i-- {
			select {
			case <-c:
				return
			default:
				doSomeThing(i)
				if i == 4 {
					lastIndex = i
					return
				}
			}
		}
	})
	if lastIndex != 4 {
		t.Errorf("lastIndex error: want=%d,actual=%d", 4, lastIndex)
	}
	time.Sleep(time.Second)
	fmt.Println("finished")
}

func TestAfterStopWithContext(t *testing.T) {
	var lastIndex int
	AfterStopWithContext(10*time.Second, func(c context.Context) {
		for i := 5; i >= 0; i-- {
			select {
			case <-c.Done():
				return
			default:
				doSomeThing(i)
				if i == 4 {
					lastIndex = i
					return
				}
			}
		}
	})
	if lastIndex != 4 {
		t.Errorf("lastIndex error: want=%d,actual=%d", 4, lastIndex)
	}
	time.Sleep(time.Second)
	fmt.Println("finished")
}

func doSomeThing(i int) {
	time.Sleep(1 * time.Second)
	fmt.Printf("do %d something...\n", i)
}
