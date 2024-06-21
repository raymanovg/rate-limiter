package fixed_window

import (
	"sync"
	"time"
)

const (
	defaultCapacity = 10
	defaultWindow   = 5 * time.Second
)

type (
	Option      func(*FixedWindow)
	FixedWindow struct {
		mu       sync.Mutex
		capacity int
		window   time.Duration
		counter  int
	}
)

func WithCapacity(capacity int) Option {
	return func(w *FixedWindow) {
		w.capacity = capacity
	}
}

func WithWindow(window time.Duration) Option {
	return func(w *FixedWindow) {
		w.window = window
	}
}

func NewFixedWindow(options ...Option) *FixedWindow {
	fw := &FixedWindow{
		mu:       sync.Mutex{},
		capacity: defaultCapacity,
		window:   defaultWindow,
		counter:  0,
	}

	for _, option := range options {
		option(fw)
	}

	go fw.StarTicker()

	return fw
}

func (fw *FixedWindow) Allow() bool {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	if fw.counter >= fw.capacity {
		return false
	}

	fw.counter++

	return true
}

func (fw *FixedWindow) StarTicker() bool {
	ticker := time.NewTicker(fw.window * time.Second)

	for {
		<-ticker.C
		fw.mu.Lock()
		fw.counter = 0
		fw.mu.Unlock()
	}
}
