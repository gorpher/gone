package cache

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value interface{}) error
	Del(key string) error

	// SetWithTTL 设置key，超时时间，注意单位，不能直接使用秒,比如一分钟必须写成 60*time.Second
	SetWithTTL(key string, value string, duration time.Duration) error
}

type Options struct {
	cacheDir     string
	redis        *RedisOptions
	redisCmdable redis.Cmdable
	inMemery     bool
}

type OptFunc func(*Options) *Options

func WithCacheDir(cacheDir string) OptFunc {
	return func(opt *Options) *Options {
		opt.cacheDir = cacheDir
		return opt
	}
}
func WithInMemory() OptFunc {
	return func(opt *Options) *Options {
		opt.inMemery = true
		return opt
	}
}

func WithRedis(redis *RedisOptions) OptFunc {
	return func(opt *Options) *Options {
		opt.redis = redis
		return opt
	}
}

func WithRedisCmdable(redisCmdable redis.Cmdable) OptFunc {
	return func(opt *Options) *Options {
		opt.redisCmdable = redisCmdable
		return opt
	}
}

var _ Cache = (*MemoryCache)(nil)

var _ Cache = (*BoltCache)(nil)

func NewCache(opts ...OptFunc) (Cache, error) {
	var options = new(Options)
	options.cacheDir = ".cache"
	for _, opt := range opts {
		opt(options)
	}

	if options.redisCmdable != nil {
		return NewRedisCacheDB(options.redisCmdable)
	}

	if options.redis != nil && len(options.redis.Address) > 0 {
		return NewRedisCache(options.redis)
	}
	if options.inMemery {
		return NewBadgerCache(options.cacheDir, true)
	}
	return NewBadgerCache(options.cacheDir, false)
}
