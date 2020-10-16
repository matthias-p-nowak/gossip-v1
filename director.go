package main

import (
  "runtime"
  "sync"
  "time"
  "regexp"
  "log"
  "strings"
)

// used as index for the different maps
const (
  Number = iota
  CallId
  Via
  DirectoryEnd
)

// maps a certain tag to a channel for gossip items
type director map[string]chan *GossipItem

var (
  // the maps from tags to strings
  DirectorChans []director
  // the related mutex for access serialization
  NumberLock []sync.RWMutex
  viaReg      *regexp.Regexp
  siptelReg   *regexp.Regexp
)

// when the goroutines serving the channels end, they might have been registered several times
func cleanUpDirector() {
  for {
    for i := 0; i < DirectoryEnd; i++ {
      time.Sleep(time.Second)
      k2d := []string{}
      NumberLock[i].RLock()
      for k, v := range DirectorChans[i] {
        select {
        case v <- nil:
          // everything ok
        default:
          // channel is filled == no reader left
          k2d = append(k2d, k)
        }
      }
      NumberLock[i].RUnlock()
      NumberLock[i].Lock()
      for _, k := range k2d {
        delete(DirectorChans[i], k)
      }
      NumberLock[i].Unlock()
      runtime.Gosched()
    }
  }
}

func RegisterChan(dir int, key string, ch chan *GossipItem) {
  NumberLock[dir].Lock()
  DirectorChans[dir][key] = ch
  NumberLock[dir].Unlock()
}

func FillChan(ch chan *GossipItem) {
  for {
    select {
    case ch <- nil:
      // ok
    default:
      return
    }
  }
}

func SendItem(dir int, key string, it *GossipItem) (ok bool) {
  NumberLock[dir].RLock()
  ch := DirectorChans[dir][key]
  NumberLock[dir].RUnlock()
  if ch == nil {
    return false
  }
  select {
  case ch <- it:
    ok = true
  default:
    ok = false
  }
  return
}

func DirectItem(item *GossipItem){
  msg:=item.msg
  hd:=msg.Header["Via"]
  if len(hd)>0 {
    str:=hd[0]
    m := viaReg.FindStringSubmatch(str)
    if m != nil && len(m) > 1 {
      via := m[1]
      if SendItem(Via,via,item) { return }
    }
  }
  hd=msg.Header["Call-ID"]
  if len(hd)>0 {
    str:=strings.TrimSpace(hd[0])
    if SendItem(CallId,str,item) { return }
  }
  m:= siptelReg.FindStringSubmatch(msg.SipLine)
  if m != nil && len(m)> 1{
    telno:=m[1]
    if SendItem(Number,telno,item) { return }
  }
  log.Fatal("don't know where to send ",item)
}

func init() {
  NumberLock = make([]sync.RWMutex, DirectoryEnd)
  DirectorChans = make([]director, DirectoryEnd)
  for i := 0; i < DirectoryEnd; i++ {
    DirectorChans[i] = make(director)
  }
  go cleanUpDirector()
  re, err := regexp.Compile("branch=([^; ]*)")
  if err != nil { log.Fatal(err) }
  viaReg = re
  re, err = regexp.Compile(":([^@; ]*)")
  if err != nil { log.Fatal(err) }
  siptelReg = re
}

