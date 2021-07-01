package gone

import (
	"context"
	"time"
)

//AfterStopWithContext 超时停止协程，关闭协程。使用context实现
func AfterStopWithContext(d time.Duration, f func(context.Context)) {
	c, cancelFunc := context.WithTimeout(context.Background(), d)
	defer cancelFunc()
	go func() {
		f(c)
		cancelFunc()
	}()
	<-c.Done()
}

//AfterStopFunc 超时停止协程，关闭协程。使用原始channel实现
//这里新建了两个channel是防止关闭一个已经关闭的channel导致panic，这里还有优化点。
//该函数主使用场景，比如防止扫描磁盘时间过长，在规定时间里获取结果。
func AfterStopFunc(d time.Duration, f func(<-chan struct{})) {
	c := make(chan struct{})
	t := make(chan struct{})
	go func() {
		f(c)
		close(t)
	}()
	select {
	case <-time.After(d):
	case <-t:
	}
	close(c)
}
