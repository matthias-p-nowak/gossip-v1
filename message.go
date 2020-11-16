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

var (
  // Alphabet comprises the symbols for creating random strings
  Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GossipMsgHeadersHeaders: Headers might have multiple values (via, record-routes, route),
// hence it contains a list of values
type GossipMsgHeaders map[string][]string

// GossipMsg contains the SIP message
type GossipMsg struct {
  // StartLine: the first line of both Request and response
  StartLine string
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
