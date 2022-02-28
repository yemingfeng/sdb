package server

import (
	"github.com/bsm/ratelimit"
	"time"
)

type RateLimit struct {
	rateLimiter *ratelimit.RateLimiter
}

func CreateRateLimit(rate int) *RateLimit {
	return &RateLimit{
		rateLimiter: ratelimit.New(rate, time.Second),
	}
}

func (rateLimit *RateLimit) Limit() bool {
	return rateLimit.rateLimiter.Limit()
}
