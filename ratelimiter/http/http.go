// Package http rate decorates an http.Handler and prevents too many requests being made
package http

import (
	"net/http"
	"strings"
	"time"
	"fmt"
)

const StatusTooManyRequests = 429

type RateLimitedHandler struct {
	http.Handler
	rateLimiter RateLimiter
}

type RateLimiter func(*http.Request) (ok bool, retryAfter time.Duration)

func Decorate(delegate http.Handler, rateLimiter RateLimiter) RateLimitedHandler {
	if rateLimiter == nil {
		rateLimiter = DefaultRateLimiter
	}
	return RateLimitedHandler{delegate, rateLimiter}
}

func (lh RateLimitedHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if ok, retryAfter := lh.rateLimiter(r); ok {
		lh.Handler.ServeHTTP(rw, r)
	} else {
		rw.Header().Set("Retry-After", fmt.Sprint(int64(retryAfter.Seconds())))
		http.Error(rw, "Too many requests", StatusTooManyRequests)
	}
}
