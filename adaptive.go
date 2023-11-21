package limit

// ## 自适应限流算法

// 参考 tcp bbr 算法、go-zero、kratos 实现

// 算法原理其实很简单：
// 就是根据 cpu 负载和请求负载来判断要不要限流，有一些要点如下：
// 1. 也可以使用其它设备负载，不过 cpu 显然是限制最大的，所以就是用当前请求的 cpu 使用率来判断即可，而 cpu 负载的统计可以每次都获取，也可以定时获取
// 2. 请求负载怎么统计？其实很简单，就是获取历史上能够通过的最大请求量，然后再根据当前的请求量来判断，要不要限流，这里展开也有几个问题：
// 2.1 只是历史请求量就行吗？其实不然，因为机器的负载一直在变动，最好取最近一段时间的最大请求量比较好，比如前 5min 的
// 2.2 怎么统计当前的请求量？这个可能需要处理一下
// 2.3 怎么统计过去 x 分钟的最大请求量？这个可以根据滑动窗口算法来采集

// 一般来说，直接根据这两个负载就能够判断了，不过 kratos 还做了优化，就是当 cpu 负载小于 80% 后，不是直接判断请求负载，而是先判断一下上一次限流的时间，然后再去判断请求负载，这样是为了减少限流之后 cpu 负载降低那瞬间出现的毛刺

// ## 参考
// https://lailin.xyz/post/go-training-week6-4-auto-limiter.html
// https://github.com/go-kratos/kratos/blob/v1.0.x/pkg/ratelimit/bbr/bbr.go
// https://github.com/sado0823/go-bbr-ratelimit

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// 自适应限流算法
// todo: 当前请求量的统计、过去 5min 最大请求量的统计（滑动窗口采集）
type limitAdaptive struct {
	lastLimitTime *time.Time // 上一次限流时间
	requestNum    int        // 当前请求量
	lastLoad      int        // 过去 5min 内允许的最大吞吐量
}

func NewAdaptiveLimit(c *Config) *limitAdaptive {
	l := &limitAdaptive{}

	return l
}

// 只是占位
func (li *limitAdaptive) Wait() {
	panic("not implemented") // TODO: Implement
}

func (li *limitAdaptive) Allow() bool {
	// [step 1] 获取 cpu 使用率
	f, err := cpu.Percent(0, false)
	if err != nil {
		log.Printf("get cpu use info failed when allow, error:%s", err)
		return true
	}

	// [step 2] 获取当前时间
	t := time.Now()

	// [step 3] 判断 cpu 负载是否小于 80%
	if f[0] < 0.8 {
		// [step 3.1] 判断上一次限流时间是否在 1s 内
		if t.Sub(*li.lastLimitTime) < time.Second {
			// [step 3.1.1] 判断请求负载是否超过,超过就直接拒绝
			if li.requestNum > li.lastLoad {
				return false
			}
		}
		// [step 3.2] 如果上一次限流在 1s 外、或者请求量没有超过负载，则允许
		return true
	}

	// [step 4] 如果 cpu 负载超过了 80%， 则再根据请求负载判断是否要限流
	return li.requestNum > li.lastLoad
}
