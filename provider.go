package main

import (
  "log"
  "strings"
  "sync"
  "net"
)

type Closer interface {
  Close() error
}

var (
  provider sync.WaitGroup
  networks=make(map[string]Closer)
  mNetworks sync.Mutex
  localEPs=make(map[string]chan *GossipItem)
  mLocalEPs sync.Mutex
  retransmit=make(map[string] *GossipItem)
  mRetransmit sync.Mutex
)

func addNetwork(provider string, cl Closer){
    mNetworks.Lock()
  networks[provider]=cl
  mNetworks.Unlock()
}

func addLocalEP(provider string,ch chan *GossipItem){
  mLocalEPs.Lock()
  localEPs[provider]=ch
  mLocalEPs.Unlock()
}

func EndProviders(){
  mNetworks.Lock()
  for _,c:=range networks{
    c.Close()
  }
  mNetworks.Unlock()
  provider.Wait()
}

type UdpGossipProvider struct {
  ch chan *GossipItem
  netConn net.PacketConn
}



func (p *UdpGossipProvider) Receiver() {
}

/*
 * receives items from the internal system and send it out on the channel, one after the other
 */
func (p *UdpGossipProvider) Sender() {
  for item := range p.ch {
    if item == nil {
      continue
    }
  }
}

func newUdpProvider(provider string)(p *UdpGossipProvider, err error) {
  p=new(UdpGossipProvider)
  p.ch=make(chan *GossipItem,8)
  parts := strings.Split(provider, "/")
  nc, err := net.ListenPacket("udp", parts[1])
  p.netConn=nc
  addNetwork(provider,nc)
  addLocalEP(provider,p.ch)
  return
}

func NewProvider(provider string) {
  parts := strings.Split(provider, "/")
  if len(parts) != 2 {
    log.Fatal("provider string should be like 'udp/192.168.0.8:5060', got " + provider)
  }
  switch parts[0] {
  case "udp":
    p,err:=newUdpProvider(provider)
    if err != nil { log.Fatal(err)}
    go p.Sender()
    go p.Receiver()
  default:
    log.Fatal("couldn't handle network " + parts[0])
  }
}
