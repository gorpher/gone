package osutil

import "time"

// SlowFunc 执行slowFunc函数，如果超过阈值会执行callback函数，不会等待sf函数
func SlowFunc(threshold time.Duration, sf func() error, cb func()) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- sf()
	}()
	select {
	case err := <-errCh:
		return err
	case <-time.After(threshold):
		cb()
		return nil
	}
}
