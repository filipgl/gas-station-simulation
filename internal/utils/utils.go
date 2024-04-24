package utils

import (
	"math/rand"
	"time"
)

func RandomIntFromRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomDurationFromRange(min, max time.Duration) time.Duration {
	mn, mx := min.Nanoseconds(), max.Nanoseconds()
	return time.Duration(rand.Int63n(mx-mn+1) + mn)
}
