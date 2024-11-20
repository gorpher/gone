package gormzerolog

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"

	"gorm.io/gorm/logger"
)

// NewLoggerInterface  level,slowThreshold,skipFrameCount
func NewLoggerInterface(isDebug bool, args ...int) logger.Interface {
	g := &gormLogger{}
	if len(args) > 0 {
		g.slowThreshold = time.Duration(args[0])
	} else {
		g.slowThreshold = 100 * time.Millisecond
	}
	if len(args) > 1 {
		g.skipFrameCount = args[1]
	}
	g.isDebug = isDebug
	return g
}

type gormLogger struct {
	isDebug        bool
	slowThreshold  time.Duration
	skipFrameCount int
}

type gormZeroLogSilentKey struct {
}

func WithContextSilent() context.Context {
	return context.WithValue(context.Background(), gormZeroLogSilentKey{}, true)
}

func (l gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l gormLogger) Error(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Error().Msg(fmt.Sprintf(msg, opts...))
}

func (l gormLogger) Warn(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Warn().Msg(fmt.Sprintf(msg, opts...))
}

func (l gormLogger) Info(ctx context.Context, msg string, opts ...interface{}) {
	zerolog.Ctx(ctx).Info().Msg(fmt.Sprintf(msg, opts...))
}

func (l gormLogger) Trace(ctx context.Context, begin time.Time, f func() (string, int64), err error) {
	zl := zerolog.Ctx(ctx)
	var event *zerolog.Event
	if silent, ok := ctx.Value(gormZeroLogSilentKey{}).(bool); !ok || !silent {
		if l.isDebug {
			event = zl.Debug()
		} else {
			event = zl.Trace()
		}
	}
	var durKey string

	switch zerolog.DurationFieldUnit {
	case time.Nanosecond:
		durKey = "elapsed_ns"
	case time.Microsecond:
		durKey = "elapsed_us"
	case time.Millisecond:
		durKey = "elapsed_ms"
	case time.Second:
		durKey = "elapsed"
	case time.Minute:
		durKey = "elapsed_min"
	case time.Hour:
		durKey = "elapsed_hr"
	default:
		durKey = "elapsed_"
	}
	elapsed := time.Since(begin)
	sql, rows := f()
	if event == nil {
		return
	}
	event.Dur(durKey, elapsed)
	if l.skipFrameCount > 0 {
		event.CallerSkipFrame(zerolog.CallerSkipFrameCount + l.skipFrameCount)
	}
	if sql != "" {
		event.Str("sql", sql)
	}
	if rows > -1 {
		event.Int64("rows", rows)
	}

	if l.slowThreshold > 0 && elapsed > l.slowThreshold {
		event.Str("slowLog", fmt.Sprintf("SLOW SQL >= %v", l.slowThreshold))
	}
	event.Send()

	return
}
