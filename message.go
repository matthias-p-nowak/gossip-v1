package main

import(
  "net"
)

type GossipMsgHeaders struct {
  To string
  From string
  CallId string
  CSeg int
  Via []string
  Route []string
  Rroute []string
  Others []string
}

type GossipMsg struct {
  Network string
  Local net.Addr
  Remote net.Addr
  Headers GossipMsgHeaders
  Body string
}
