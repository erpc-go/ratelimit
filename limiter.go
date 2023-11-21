package limit

import (
	"time"
)

// Limiter
type Limiter interface {
	Wait()       // 同步睡眠
	Allow() bool // 异步返回 bool
}

func New(rate int64, opts ...Option) (l Limiter) {
	c := &Config{
		Rate:      rate,
		Circle:    time.Second, // 默认一秒钟
		LimitType: TypeWindows, // 默认 xx 算法实现
	}

	for _, opt := range opts {
		opt(c)
	}

	switch c.LimitType {
	case TypeWindows:
		l = NewWindowLimit(c)
	case TypeCount:
		l = NewCountLimit(c)
	case TypeTokenBucket:
		l = NewTokenBucketLimit(c)
	case TypeLeakyBucket:
		l = NewLeakyBucketLimit(c)
	case TypeLimitChannel:
		l = NewChannelLimit(c)
	case TypeAdaptive:
		l = NewAdaptiveLimit(c)
	default:
		l = NewTokenBucketLimit(c)
	}

	return
}
