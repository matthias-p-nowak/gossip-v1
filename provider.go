package main

import (
  "bufio"
  "log"
  "net"
  "strconv"
  "strings"
  "sync"
  "hash/adler32"
  "fmt"
  "bytes"
)

type Closer interface {
  Close() error
}

const (
  BufSize = 131072
)

var (
  provider    sync.WaitGroup
  networks    = make(map[string]Closer)
  mNetworks   sync.Mutex
  localEPs    = make(map[string]chan *GossipItem)
  mLocalEPs   sync.Mutex
  retransmit  = make(map[string]*GossipItem)
  mRetransmit sync.Mutex
  viaMsgs     = make(map[string]*GossipMsg)
  mViaMsgs    sync.Mutex
)


func addNetwork(provider string, cl Closer) {
  mNetworks.Lock()
  networks[provider] = cl
  mNetworks.Unlock()
}

func addLocalEP(provider string, ch chan *GossipItem) {
  mLocalEPs.Lock()
  localEPs[provider] = ch
  mLocalEPs.Unlock()
}

func EndProviders() {
  mNetworks.Lock()
  for _, c := range networks {
    c.Close()
  }
  mNetworks.Unlock()
  provider.Wait()
}

func ScanPost(sc *bufio.Scanner, laddr, raddr net.Addr,ch chan *GossipItem) (err error) {
  msg := new(GossipMsg)
  msg.Direction = MsgIn
  hash:=adler32.New()
  // first line
  msg.SipLine = sc.Text()
  hash.Write([]byte(msg.SipLine))
  for sc.Scan() {
    str := sc.Text()
    hash.Write([]byte(str))
    parts := strings.SplitN(str, ":", 2)
    if len(parts) == 2 {
      msg.Header[parts[0]] = append(msg.Header[parts[0]], parts[1])
    }
  }
  cla := msg.Header["Content-Length"]
  l := 0
  if len(cla) > 0 {
    cl := strings.TrimSpace(cla[0])
    l, err = strconv.Atoi(cl)
    if err != nil {
      return
    }
  }
  got := 0
  for got < l && sc.Scan() {
    str := sc.Text()
    hash.Write([]byte(str))
    got += len(str) + 2
    msg.Body = append(msg.Body, str)
  }
  vh := msg.Header["Via"]
  if vh != nil {
    m := viaReg.FindStringSubmatch(vh[0])
    if m != nil && len(m) > 1 {
      via := m[1]
      mViaMsgs.Lock()
      msg2 := viaMsgs[via]
      msg2.RetrCount = NoRetrans
      delete(viaMsgs, via)
      mViaMsgs.Unlock()
    }
  }
  item:=new(GossipItem)
  item.msg=msg
  item.Hash=hash.Sum32()
  if cfg.Verbose >= VerboseMessages {
    fmt.Printf("Msg: %s\n",msg.SipLine)
  }
  DirectItem(item)
  return
}

type UdpGossipProvider struct {
  ch      chan *GossipItem
  netConn net.PacketConn
}

type TcpGossipProvider struct {
  ch      chan *GossipItem
  netConn *net.TCPListener
  conns   map[string]*net.TCPConn
  mConns  sync.Mutex
}

func (p *UdpGossipProvider) Receiver() {
  // TODO: implement reading...
  lAddr:=p.netConn.LocalAddr()
  buf := make([]byte, BufSize)
  for{
    n,rAddr,err:=p.netConn.ReadFrom(buf)
    if err!=nil { log.Fatal(err)}
    bb:=buf[:n]
    scan:=bufio.NewScanner(bytes.NewReader(bb))
     err = ScanPost(scan,lAddr, rAddr, p.ch)
    if err != nil {
      log.Fatal(err)
    }
  }
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

func newUdpProvider(provider string) (p *UdpGossipProvider, err error) {
  p = new(UdpGossipProvider)
  p.ch = make(chan *GossipItem, 8)
  parts := strings.Split(provider, "/")
  nc, err := net.ListenPacket("udp", parts[1])
  p.netConn = nc
  addNetwork(provider, nc)
  addLocalEP(provider, p.ch)
  return
}

func newTcpProvider(provider string) (p *TcpGossipProvider, err error) {
  p = new(TcpGossipProvider)
  p.ch = make(chan *GossipItem, 8)
  p.conns = make(map[string]*net.TCPConn)
  parts := strings.Split(provider, "/")
  addr, err := net.ResolveTCPAddr("tcp", parts[1])
  if err != nil {
    log.Fatal(err)
  }
  nc, err := net.ListenTCP("tcp", addr)
  p.netConn = nc
  return
}

func (p *TcpGossipProvider) Sender() {
  for item := range p.ch {
    if item == nil {
      continue
    }
  }
}

func (p *TcpGossipProvider) ReceiveStream(conn *net.TCPConn) {
  laddr := conn.LocalAddr()
  raddr := conn.RemoteAddr()
  hp := raddr.String()
  p.mConns.Lock()
  p.conns[hp] = conn
  p.mConns.Unlock()
  scanner := bufio.NewScanner(conn)
  for scanner.Scan() {
    err := ScanPost(scanner,laddr, raddr, p.ch)
    if err != nil {
      log.Fatal(err)
    }
  }
  p.mConns.Lock()
  delete(p.conns, hp)
  p.mConns.Unlock()
}

func (p *TcpGossipProvider) Receiver() {
  for {
    conn, err := p.netConn.AcceptTCP()
    if err != nil {
      log.Fatal(err)
    }
    go p.ReceiveStream(conn)
  }
}

func NewProvider(provider string) {
  parts := strings.Split(provider, "/")
  if len(parts) != 2 {
    log.Fatal("provider string should be like 'udp/192.168.0.8:5060', got " + provider)
  }
  switch parts[0] {
  case "udp":
    p, err := newUdpProvider(provider)
    if err != nil {
      log.Fatal(err)
    }
    go p.Sender()
    go p.Receiver()
  case "tcp":
    p, err := newTcpProvider(provider)
    if err != nil {
      log.Fatal(err)
    }
    go p.Sender()
    go p.Receiver()
  default:
    log.Fatal("couldn't handle network " + parts[0])
  }
}
