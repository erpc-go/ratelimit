package limit

import (
	"fmt"
	"testing"
	"time"

	"github.com/erpc-go/limit/tools"
)

func TestNewTokenBucketLimit(t *testing.T) {
}

func TestLimitTokenBucket_Wait(t *testing.T) {
	l := NewTokenBucketLimit(&Config{
		Rate:   1000,
		Circle: time.Second,
	})

	pre := time.Now()
	for i := 0; i < 10000; i++ {
		l.Wait()
		t := time.Now()
		fmt.Println(t.Sub(pre))
		pre = t
	}
}

func TestLimitTokenBucket_Allow(t *testing.T) {
	l := NewTokenBucketLimit(&Config{
		Rate:   1000,
		Circle: time.Second,
	})

	data := make([]tools.Item, 0)

	for i := 0; i < 1000000; i++ {
		fmt.Println(l.Allow())
		// data = append(data, tools.Item{
		// 	Time: time.Now(),
		// 	Data: l.Allow(),
		// })
	}

	tools.Listen(8899, data)
}
