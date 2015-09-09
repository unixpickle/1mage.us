package main

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

var rateLimitHostTable map[string]int = map[string]int{}
var rateLimitLastTime time.Time = time.Now()
var rateLimitLock sync.Mutex

// RateLimitRequest should be called once per uploaded image.
// This will return true if the requester has reached its rate limit.
func RateLimitRequest(r *http.Request) bool {
	// TODO: in the future, there will be some sort of cookie to bypass rate limiting.

	rateLimitLock.Lock()
	defer rateLimitLock.Unlock()

	now := time.Now()
	if rateLimitLastTime.Add(time.Hour).Before(now) || rateLimitLastTime.Hour() != now.Hour() {
		rateLimitHostTable = map[string]int{}
	}

	host := r.RemoteAddr
	if forwardHeader := r.Header.Get("X-Forwarded-For"); forwardHeader != "" {
		forwardHosts := strings.Split(forwardHeader, ", ")
		host = forwardHosts[0]
	}

	if count, ok := rateLimitHostTable[host]; !ok {
		rateLimitHostTable[host] = 1
		return 1 <= rateLimitMaxCount()
	} else {
		rateLimitHostTable[host] = count + 1
		return count < rateLimitMaxCount()
	}
}

func rateLimitMaxCount() int {
	GlobalDatabase.RLock()
	defer GlobalDatabase.RUnlock()
	return GlobalDatabase.Config.MaxCountPerHour
}
