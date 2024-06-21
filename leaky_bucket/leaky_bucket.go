package leaky_bucket

import "time"

type LeakyBuket struct {
	bucket   chan chan struct{}
	capacity int
	rate     float64
}

const (
	defaultCapacity = 100
	defaultRate     = 10
)

type Option func(*LeakyBuket)

func WithCapacity(capacity int) Option {
	return func(lb *LeakyBuket) {
		lb.capacity = capacity
	}
}

func WithRate(rate float64) Option {
	return func(lb *LeakyBuket) {
		lb.rate = rate
	}
}

func NewLeakyBucket(opts ...Option) *LeakyBuket {
	lb := &LeakyBuket{
		bucket:   make(chan chan struct{}, defaultCapacity),
		capacity: defaultCapacity,
		rate:     defaultRate,
	}

	for _, opt := range opts {
		opt(lb)
	}

	go lb.leak()

	return lb
}

func (lb *LeakyBuket) Allow() bool {
	if len(lb.bucket) >= lb.capacity {
		return false
	}

	ch := make(chan struct{})
	lb.bucket <- ch
	<-ch

	return true
}

func (lb *LeakyBuket) leak() {
	for ch := range lb.bucket {
		ch <- struct{}{}
		rpms := int64(1000000 / lb.rate)
		rpms = max(1, rpms)

		time.Sleep(time.Duration(rpms) * time.Microsecond)
	}
}
