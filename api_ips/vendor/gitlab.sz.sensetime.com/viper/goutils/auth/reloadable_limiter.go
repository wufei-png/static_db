package auth

import (
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
)

type ReloadableLimiter struct {
	oldRate int64
	l       atomic.Value
}

func createBucket(rate int64) *ratelimit.Bucket {
	interval := time.Second / time.Duration(rate)
	if interval <= 0 {
		interval = time.Nanosecond
	}

	return ratelimit.NewBucket(interval, rate)
}

func NewReloadableLimiter(rate int64) *ReloadableLimiter {
	if rate <= 0 {
		panic("rate should > 0")
	}
	r := &ReloadableLimiter{
		oldRate: rate,
	}
	r.l.Store(createBucket(rate))
	return r
}

func (r *ReloadableLimiter) Reset(rate int64) bool {
	if rate <= 0 {
		panic("rate should > 0")
	}
	old := atomic.SwapInt64(&r.oldRate, rate)
	if old == rate {
		return false
	}
	r.l.Store(createBucket(rate))
	return true
}

func (r *ReloadableLimiter) getBucket() *ratelimit.Bucket {
	return r.l.Load().(*ratelimit.Bucket)
}

func (r *ReloadableLimiter) Wait(count int64) {
	r.getBucket().Wait(count)
}

func (r *ReloadableLimiter) WaitMaxDuration(count int64, maxWait time.Duration) bool {
	return r.getBucket().WaitMaxDuration(count, maxWait)
}

func (r *ReloadableLimiter) Take(count int64) time.Duration {
	return r.getBucket().Take(count)
}

func (r *ReloadableLimiter) TakeMaxDuration(count int64, maxWait time.Duration) (time.Duration, bool) {
	return r.getBucket().TakeMaxDuration(count, maxWait)
}

func (r *ReloadableLimiter) TakeAvailable(count int64) int64 {
	return r.getBucket().TakeAvailable(count)
}

func (r *ReloadableLimiter) Available() int64 {
	return r.getBucket().Available()
}

func (r *ReloadableLimiter) Capacity() int64 {
	return r.getBucket().Capacity()
}

func (r *ReloadableLimiter) Rate() float64 {
	return r.getBucket().Rate()
}
