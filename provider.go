package main

import(
  "log"
  "strings"
)

type UdpGossipProvider struct{
  ch chan *GossipItem

}

func (g *UdpGossipProvider) Receiver(){
}

/*
 * receives items from the internal system and send it out on the channel, one after the other
 */
func (g *UdpGossipProvider) Sender(){
  for item:=range(g.ch){
    if item==nil { continue }
  }
}

func NewProvider(provider string){
  parts:=strings.Split(provider,"/")
  if len(parts)!=2 {
    log.Fatal("provider string should be like 'udp/localhost:5060', got "+provider)
  }
  switch parts[0]{
    case "udp":

    default:
      log.Fatal("couldn't handle network "+parts[0])
  }
}
