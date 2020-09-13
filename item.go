package main

import (
  "time"
)

//
type GossipItem struct {
  msg *GossipMsg
}

func delaySend(dur time.Duration, ch chan *GossipItem, gi *GossipItem){
  time.AfterFunc(dur, func(){
    select{
      case ch <- gi:
      default:
    }
  })
}
