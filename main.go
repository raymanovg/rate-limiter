package main

import (
	"fmt"
	"github.com/raymanovg/rate_limiter/fixed_window"
	"github.com/raymanovg/rate_limiter/leaky_bucket"
	"github.com/raymanovg/rate_limiter/sliding_window"
	"net/http"

	"github.com/raymanovg/rate_limiter/token_bucket"
)

type Limiter interface {
	Allow() bool
}

func main() {
	limiters := map[string]Limiter{
		"token_bucket": token_bucket.NewTokenBucket(
			token_bucket.WithBurst(5),
			token_bucket.WithRate(5),
		),
		"leaky_bucket": leaky_bucket.NewLeakyBucket(
			leaky_bucket.WithCapacity(10),
			leaky_bucket.WithRate(2.0),
		),
		"fixed_window": fixed_window.NewFixedWindow(
			fixed_window.WithCapacity(10),
			fixed_window.WithWindow(5),
		),
		"sliding_window": sliding_window.NewSlidingWindow(
			sliding_window.WithWindow(5),
			sliding_window.WithRate(2.0),
		),
	}

	for name, limiter := range limiters {
		http.HandleFunc(fmt.Sprintf("/%s", name), func(w http.ResponseWriter, r *http.Request) {
			if limiter.Allow() {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			} else {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("too many requests"))
			}
		})
	}

	http.ListenAndServe(":8000", nil)
}
