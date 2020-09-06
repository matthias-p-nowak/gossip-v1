package main

import (
  "time"
)

type GossipItem struct {
  expired bool
  done bool
}

func delaySend(ms int , ch chan *GossipItem, gi *GossipItem){
  time.AfterFunc(time.Duration(ms*1000), func(){
    ch <- gi
  })
}
