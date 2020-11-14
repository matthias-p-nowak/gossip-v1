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
	MsgIncoming
	MsgOutgoing
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
	// Alphabet comprises the symbols for creating random strings
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GossipMsgHeadersHeaders: Headers might have multiple values (via, record-routes, route),
// hence it contains a list of values
type GossipMsgHeaders map[string][]string

// GossipMsg contains the SIP message
type GossipMsg struct {
	// SipLine: the first line of both Request and response
	SipLine string
	// normal Headers
	Header GossipMsgHeaders
	// extra Parameters for constructing the first Route-Header
	RouteParams map[string]string
	// Body line by line
	Body []string
	// RetrCount: retransmission count
	RetrCount int
	// Direction: outgoing or incoming
	Direction int
}

// GossipTransaction: details
type GossipTransaction struct {
	State     int
	ViaBranch string
	Dialog    *GossipDialog
	Pos       int
}

// GossipDialog: dialog details
// To/From contains the Remote/Local tag dependend on direction and Request/Response
type GossipDialog struct {
	LocalTag      string
	RemoteTag     string
	RemoteUrl     string // who we are calling/called by
	RemoteContact string
	Transactions  []*GossipTransaction
	Call          *GossipCall
	Pos           int
}

// GossipCall: details
type GossipCall struct {
	CallId       string // constant
	CallSeq      int    // increasing
	Status       int    `CallStatus`
	LocalContact string
	Dialogs      []*GossipDialog
}

func init() {
	t := time.Now().Unix()
	rand.Seed(t)
}

// RandStrings returns a string of length <l>
func RandString(l int) string {
	aLen := len(Alphabet)
	bb := make([]byte, l)
	for i := 0; i < l; i++ {
		bb[i] = Alphabet[rand.Intn(aLen)]
	}
	return string(bb)
}

// NewTransaction
func (gd *GossipDialog) NewTransaction() (nt *GossipTransaction) {
	nt = new(GossipTransaction)
	nt.Dialog = gd
	nt.Pos = len(gd.Transactions)
	gd.Transactions = append(gd.Transactions, nt)
	return
}

// NewDialog
func (gc *GossipCall) NewDialog() (nd *GossipDialog) {
	nd = new(GossipDialog)
	nd.Call = gc
	nd.Pos = len(gc.Dialogs)
	gc.Dialogs = append(gc.Dialogs, nd)
	return
}
