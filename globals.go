package main

import (
	"sync"
	"time"

	"github.com/unixpickle/ratelimit"
)

// ShutdownLock should be locked for reading whenever the database is being modified.
// When the app is shutting down, this will be locked for writing to block further changes.
var ShutdownLock sync.RWMutex

// TemporaryDirectory can be used in tasks like uploading files. It ensures that the file will be
// deleted if the task terminates.
var TemporaryDirectory string

// GlobalDb is used throughout the server to manage and access data.
var GlobalDb *Db

// RateLimiter is used to rate limit users.
var RateLimiter *ratelimit.TimeSliceLimiter = ratelimit.NewTimeSliceLimiter(time.Hour, 0)

// HTTPNamer is used in the rate-limiting process.
var HTTPNamer ratelimit.HTTPRemoteNamer
