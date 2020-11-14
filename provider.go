package main

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/adler32"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Closer: need a means to shut down internet TCP connections
type Closer interface {
	Close() error
}

const (
	BufSize = 131072 // 128 KByte
)

var (
	// Maps need a lock if used concurrently
	provider  sync.WaitGroup
	networks  = make(map[string]Closer)
	mNetworks sync.Mutex
	//
	localEPs  = make(map[string]chan *GossipItem)
	mLocalEPs sync.Mutex
	//
	retransmit  = make(map[string]*GossipItem)
	mRetransmit sync.Mutex
	//
	viaMsgs  = make(map[string]*GossipMsg)
	mViaMsgs sync.Mutex
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

// EndProviders: shuts down all connections
func EndProviders() {
	mNetworks.Lock()
	for _, c := range networks {
		c.Close()
	}
	mNetworks.Unlock()
	provider.Wait()
}

// ScanPost reads one post from the scanner and gives it to the director
func ScanPost(sc *bufio.Scanner, lAddr, rAddr net.Addr, ch chan *GossipItem) (err error) {
	msg := new(GossipMsg)
	msg.Direction = MsgIncoming
	hash := adler32.New()
	// first line is already scanned
	msg.SipLine = sc.Text()
	hash.Write([]byte(msg.SipLine))
	for sc.Scan() {
		str := sc.Text()
		hash.Write([]byte(str))
		parts := strings.SplitN(str, ":", 2)
		if len(parts) == 2 {
			msg.Header[parts[0]] = append(msg.Header[parts[0]], parts[1])
		} else {
			// empty line without a header
			break
		}
	}
	// find out how much to read from the scanner
	cla := msg.Header["Content-Length"]
	l := 0
	if len(cla) > 0 {
		cl := strings.TrimSpace(cla[0])
		l, err = strconv.Atoi(cl)
		if err != nil {
			return
		}
	}
	// start reading
	got := 0
	// only read the requested amount
	for got < l && sc.Scan() {
		str := sc.Text()
		hash.Write([]byte(str))
		got += len(str) + 2
		msg.Body = append(msg.Body, str)
	}
	// stopping retransmissions
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
	// preparing for sending it to the director
	item := new(GossipItem)
	item.msg = msg
	item.Hash = hash.Sum32()
	item.ch = ch
	item.localEP = lAddr.Network() + "/" + lAddr.String()
	item.remoteEP = rAddr.Network() + "/" + rAddr.String()
	if cfg.Verbose >= VerboseMessages {
		fmt.Printf("Msg: %s\n", msg.SipLine)
	}
	// the director will handle it and send it to the right analyser
	DirectItem(item)
	return
}

// UdpGossipProvider combines a network interface with a function that sends items from the channel to the net
type UdpGossipProvider struct {
	ch      chan *GossipItem // for items to send over this port
	netConn net.PacketConn
}

// TcpGossipProvider combines a listener with multiple active connections
type TcpGossipProvider struct {
	ch      chan *GossipItem // for items to be send by this provider
	netConn *net.TCPListener // listener for incoming connections
	//
	conns  map[string]*net.TCPConn // map of alive connections
	mConns sync.Mutex              // synchronizes access to conns
}

// Receiver for UDP, all messages come in the same way
func (p *UdpGossipProvider) Receiver() {
	// local side is fixed
	lAddr := p.netConn.LocalAddr()
	buf := make([]byte, BufSize)
	for {
		// stops when channel is closed
		n, rAddr, err := p.netConn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		bb := buf[:n] // need a new slice
		scan := bufio.NewScanner(bytes.NewReader(bb))
		// send it to the general message reading routine
		scan.Scan() // scanning the first line, ScanPost starts with Text()
		err = ScanPost(scan, lAddr, rAddr, p.ch)
		if err != nil {
			log.Fatal(err)
		}
	}
	// the local end is closed
}

// Sender: receives items from the internal system and send it out on the channel, one after the other
func (p *UdpGossipProvider) Sender() {
	for item := range p.ch {
		if item == nil {
			continue
		}
		// TODO: implement sending it over UDP
	}
}

// newUdpProvider creates a new UDP provider (only structure)
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

// newTcpProvider creates a new TCP provider (only structure)
func newTcpProvider(provider string) (p *TcpGossipProvider, err error) {
	p = new(TcpGossipProvider)
	p.ch = make(chan *GossipItem, 8)
	p.conns = make(map[string]*net.TCPConn)
	parts := strings.Split(provider, "/")
	addr, err := net.ResolveTCPAddr("tcp", parts[1])
	if err != nil {
		log.Fatal(err)
	}
	// creating the listener
	nc, err := net.ListenTCP("tcp", addr)
	p.netConn = nc
	return
}

// Sender: receives items on the channel and sends it over the net
func (p *TcpGossipProvider) Sender() {
	for item := range p.ch {
		if item == nil {
			continue
		}
		// TODO: implement sending on this provider
		// need to find out, if there is an existing connection - then using this one
		// otherwise, creating a new connection
	}
}

// ReceiveStream takes messages from connected streams
// streams are both from acceptance by the listener, or actively established
func (p *TcpGossipProvider) ReceiveStream(conn *net.TCPConn) {
	laddr := conn.LocalAddr()
	raddr := conn.RemoteAddr()
	hp := raddr.String()
	p.mConns.Lock()
	p.conns[hp] = conn
	p.mConns.Unlock()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		// will fail when stream is closed
		err := ScanPost(scanner, laddr, raddr, p.ch)
		if err != nil {
			log.Fatal(err)
		}
	}
	p.mConns.Lock()
	delete(p.conns, hp)
	p.mConns.Unlock()
}

// Receiver accepts new connections on the listener
func (p *TcpGossipProvider) Receiver() {
	for {
		conn, err := p.netConn.AcceptTCP()
		if err != nil {
			log.Println(err)
		}
		// start a new goroutine for receiving messages
		go p.ReceiveStream(conn)
	}
}

// NewProvider creates provider structures and starts goroutine
func NewProvider(provider string) {
	// provider is like described below
	parts := strings.Split(provider, "/")
	if len(parts) != 2 {
		log.Fatal("provider string should be like 'udp/192.168.0.8:5060', got " + provider)
	}
	// udp or tcp
	switch parts[0] {
	// UDP
	case "udp":
		p, err := newUdpProvider(provider)
		if err != nil {
			log.Fatal(err)
		}
		go p.Sender()
		go p.Receiver()
		// TCP
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
