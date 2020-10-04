package main

import (
// "net"
)

// Retransmission count
const (
	NoRetrans = iota
	ReTrOnce
	ReTrFirst
	ReTrSecond
	ReTrThird
	ReTrFourth
	ReTrFifth
	ReTrSixth
	ReTrSeventh
	ReTrEnd
)

// the message direction
const (
	Undefined = iota
	MsgIn
	MsgOut
)

// call/request status, repeating or not is on msg
const (
	SipNone         = iota
	SipInviting     // sending out invites, or 100trying on each Invite
	SipTrying       // invite send out, or invite received
	SipRinging      //
	SipEarly        // after prack
	SipEstablished  // after 200
	SipAcknowledged // now bye will end
	SipCanceling    // before final response
	SipEnding
	SipFinished
)

// transaction states
const (
	TransReqested = iota
	TransInitiated
	TransResponded
)

// session states
const (
	SipSessionOffered = iota
	SipSessionAnswered
	SipSessionEstablished
)

// Headers might have multiple values (via, record-routes, route)
type GossipMsgHeaders map[string][]string

// The SIP message with it's components
type GossipMsg struct {
	SipLine   string // both Request and response
	Header    GossipMsgHeaders
	Body      string
	RetrCount int
	Direction int
	RawMsg    []byte
}

type GossipSession struct {
	State int
}

type GossipTransaction struct {
	State     int
	ViaBranch string
}
type GossipDialog struct {
	LocalTag      string
	RemoteTag     string
	RemoteUrl     string // who we are calling/called by
	RemoteContact string
	Transactions  []GossipTransaction
}

//
type GossipCall struct {
	CallId       string
	CallSeq      int
	Status       int `CallStatus`
	LocalContact string
	Dialogs      []GossipDialog
}
