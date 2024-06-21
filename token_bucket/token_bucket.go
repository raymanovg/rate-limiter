package token_bucket

import (
	"sync"
	"time"
)

const (
	defaultBurst = 100
	defaultRate  = 1.0
)

type (
	Option      func(*TokenBucket)
	TokenBucket struct {
		mu         sync.Mutex
		tokens     int
		burst      int
		rate       float64
		lastRefill time.Time
	}
)

func WithBurst(burst int) Option {
	return func(t *TokenBucket) {
		t.burst = burst
	}
}

func WithRate(rate float64) Option {
	return func(t *TokenBucket) {
		t.rate = rate
	}
}

func NewTokenBucket(options ...Option) *TokenBucket {
	tb := &TokenBucket{
		mu:     sync.Mutex{},
		burst:  defaultBurst,
		rate:   defaultRate,
		tokens: 0,
	}

	for _, option := range options {
		option(tb)
	}

	return tb
}

func (tb *TokenBucket) Allow() bool {
	return tb.allowN(1)
}

func (tb *TokenBucket) allowN(tokens int) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens < tokens {
		return false
	}

	tb.tokens -= tokens
	return true
}

func (tb *TokenBucket) refill() {
	if tb.burst == tb.tokens {
		return
	}

	now := time.Now()
	passedSecs := now.Sub(tb.lastRefill).Seconds()
	tokensToAdd := int(tb.rate * passedSecs)
	if tokensToAdd == 0 {
		return
	}
	tb.tokens = min(tb.burst, tb.tokens+tokensToAdd)
	tb.lastRefill = now
}
