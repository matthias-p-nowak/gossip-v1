package main

import (
	"time"
)

//
type GossipItem struct {
	msg *GossipMsg
}

// for retransmissions, the same message must be send again after a certain interval
func delaySend(dur time.Duration, ch chan *GossipItem, gi *GossipItem) {
	time.AfterFunc(dur, func() {
		select {
		case ch <- gi:
		default:
			// channel is full - most likely abandoned
		}
	})
}
