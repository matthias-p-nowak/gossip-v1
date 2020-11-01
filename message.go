package main

import (
	"math/rand"
	"time"
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
	TransRequested = iota
	TransInitiated
	TransResponded
)

// session states
const (
	SipSessionOffered = iota
	SipSessionAnswered
	SipSessionEstablished
)

var (
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Headers might have multiple values (via, record-routes, route)
type GossipMsgHeaders map[string][]string

// The SIP message with it's components
type GossipMsg struct {
	SipLine   string // both Request and response
	Header    GossipMsgHeaders
	Body      []string
	RetrCount int
	Direction int
}

type GossipSession struct {
	State int
}

type GossipTransaction struct {
	State     int
	ViaBranch string
	Dialog    *GossipDialog
	Pos       int
}
type GossipDialog struct {
	LocalTag      string
	RemoteTag     string
	RemoteUrl     string // who we are calling/called by
	RemoteContact string
	Transactions  []*GossipTransaction
	Call          *GossipCall
	Pos           int
}

//
type GossipCall struct {
	CallId       string
	CallSeq      int
	Status       int `CallStatus`
	LocalContact string
	Dialogs      []*GossipDialog
}

func init() {
	t := time.Now().Unix()
	rand.Seed(t)
}

func RandString(l int) string {
	aLen := len(Alphabet)
	bb := make([]byte, l)
	for i := 0; i < l; i++ {
		bb[i] = Alphabet[rand.Intn(aLen)]
	}
	return string(bb)
}

func (gd *GossipDialog) NewTransaction() (nt *GossipTransaction) {
	nt = new(GossipTransaction)
	nt.Dialog = gd
	nt.Pos = len(gd.Transactions)
	gd.Transactions = append(gd.Transactions, nt)
	return
}

func (gc *GossipCall) NewDialog() (nd *GossipDialog) {
	nd = new(GossipDialog)
	nd.Call = gc
	nd.Pos = len(gc.Dialogs)
	gc.Dialogs = append(gc.Dialogs, nd)
	return
}
