package limit

// 计数器、固定窗口限流算法

// 计数器限流缺点：
//  1. 在两周期中间打入，即前一周期的后半部分打满和后周期的前半部分打满，会导致总体上一个周期流量达到限制的两倍
//  2. 在周期刚开始就打满，会导致该周期后面全部请求被拒绝，造成毛刺现象，请求不均匀

import (
	"time"
)

// type
type LimitType int

const (
	TypeWindows      LimitType = iota // 滑动窗口
	TypeCount                         // 计数
	TypeTokenBucket                   // 令牌桶
	TypeLeakyBucket                   // 漏桶
	TypeLimitChannel                  // 巧妙利用 channel 和定时器进行限流
	TypeAdaptive                      // 自适应算法
)

// config
type Config struct {
	Rate      int64         // 请求最多次数
	Circle    time.Duration // 请求周期
	LimitType LimitType     // 实现算法类型,默认 todo
	take      func()
}

type Option func(*Config)

func WithRate(rate int64) Option {
	return func(c *Config) {
		c.Rate = rate
	}
}

func WithLimitType(t LimitType) Option {
	return func(c *Config) {
		c.LimitType = t
	}
}

func WithCircle(t time.Duration) Option {
	return func(c *Config) {
		c.Circle = t
	}
}

func WithFunc(take func()) Option {
	return func(c *Config) {
		c.take = take
	}
}
