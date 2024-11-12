package osutil

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithContextFunc returns a copy of parent context that is cancelled when
// an os interrupt signal is received. The callback function f is invoked
// before cancellation.
func WithContextFunc(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			f()
			cancel()
		}
	}()

	return ctx
}
