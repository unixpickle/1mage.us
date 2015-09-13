package main

import (
	"net/http"
	"time"

	"github.com/unixpickle/ratelimit"
)

var rateLimiter *ratelimit.TimeSliceLimiter = ratelimit.NewTimeSliceLimiter(time.Hour, 0)
var httpNamer ratelimit.HTTPRemoteNamer

// RateLimitRequest should be called once per uploaded image.
// This will return true if the requester has reached its rate limit.
func RateLimitRequest(r *http.Request) bool {
	// TODO: in the future, there will be some sort of cookie to bypass rate limiting.
	id := httpNamer.Name(r)
	return -rateLimiter.Decrement(id) <= rateLimitMaxCount()
}

func rateLimitMaxCount() int64 {
	GlobalDatabase.RLock()
	defer GlobalDatabase.RUnlock()
	return int64(GlobalDatabase.Config.MaxCountPerHour)
}
