package main

import(
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

// Headers might have multiple values
type GossipMsgHeaders map[string][]string

// The SIP message with it's components
type GossipMsg struct {
  Url string
  Header GossipMsgHeaders
  Body string
  RetrCount int
  Direction int
  RawMsg []byte
}

type GossipCall struct {
  CallId string
  CallSeq int
}

type GossipMsgData struct{
  Call *GossipCall
}
