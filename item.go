package main

import (
	"time"
)

// GossipItem meta structure for the received SipMessage with metadata
type GossipItem struct {
	// the sip-message
	msg      *GossipMsg
	localEP  string
	remoteEP string
	// channel that answers should be send over
	ch chan *GossipItem
	// the raw packet send over IP
	RawMsg []byte
	// hash of the raw packet - for identifying retransmissions
	Hash uint32
}

// delaySend for retransmissions, the same message must be send again after a certain interval
func delaySend(dur time.Duration, ch chan *GossipItem, gi *GossipItem) {
	time.AfterFunc(dur, func() {
		select {
		case ch <- gi:
		default:
			// channel is full - most likely abandoned
		}
	})
}
