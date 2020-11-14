package main

import (
	"sync"
	"time"
)

// local data for the limiter
var (
	// ticker sends time msg in certain intervals
	ticker = time.NewTicker(time.Second)
	// number of concurrent calls - compared with cfg.Concurrent
	current int32
	// Lock and Signal
	cond *sync.Cond = sync.NewCond(new(sync.Mutex))
)

// LimiterInit: can be called several times
func LimiterInit() {
	if cfg.Rate > 0 {
		di := int64(time.Second) / int64(cfg.Rate)
		d := time.Duration(di)
		ticker = time.NewTicker(d)
	}
	// TODO: do we need to stop the old one, or does it just disappear
	if cfg.Concurrent < 1 {
		cfg.Concurrent = 1
	}
}

// FetchLimited: when requested, returns at the required rate or limits the amount of simultaneous calls
func FetchLimited() {
	cond.L.Lock()
	// no one can change current right now
	for current >= cfg.Concurrent {
		cond.Wait() // wait unlocks, waits and locks
	}
	current++
	<-ticker.C
	cond.L.Unlock()
}

// ReleaseLimited: released by the calling part
func ReleaseLimited() {
	cond.L.Lock()
	current--
	cond.Broadcast() // FetchLimited's wait can continue
	cond.L.Unlock()
}
