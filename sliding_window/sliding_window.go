package sliding_window

import (
	"sync"
	"time"
)

const (
	defaultWindow = 5 * time.Second
	defaultRate   = 2
)

type (
	Option        func(w *SlidingWindow)
	SlidingWindow struct {
		mu      sync.Mutex
		window  time.Duration
		count   int
		history []int64
	}
)

func WithWindow(window time.Duration) Option {
	return func(w *SlidingWindow) {
		w.window = window
	}
}

func WithRate(rate float64) Option {
	return func(w *SlidingWindow) {
		w.count = int(rate * float64(w.window) / float64(time.Second))
	}
}

func NewSlidingWindow(options ...Option) *SlidingWindow {
	sw := &SlidingWindow{
		mu:      sync.Mutex{},
		window:  defaultWindow,
		count:   int(defaultRate * float64(defaultWindow) / float64(time.Second)),
		history: make([]int64, 0),
	}

	for _, option := range options {
		option(sw)
	}

	return sw
}

func (sw *SlidingWindow) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now()
	for len(sw.history) > 0 && now.Sub(time.Unix(0, sw.history[0])) >= sw.window {
		sw.history = sw.history[1:]
	}

	if len(sw.history) >= sw.count {
		return false
	}

	sw.count++
	sw.history = append(sw.history, now.Unix())

	return true
}
