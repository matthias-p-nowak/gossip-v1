package main

/*
 * The directory gets items from the SIP providers and direct them to the correct tester
 */

import (
  "log"
  "regexp"
  "runtime"
  "strings"
  "sync"
  "time"
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
  viaReg     *regexp.Regexp
  siptelReg  *regexp.Regexp
)

// cleanUpDirector removes idle channels from the maps,
// idle means that the channel is full,
// when the goroutines serving the channels end, they might have been registered several times
func cleanUpDirector() {
  for gossipRunning {
    // until end of time
    for i := 0; i < DirectoryEnd; i++ {
      time.Sleep(time.Second)
      // the keys to remove
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

// RegisterChan creates entries in the director map
// dir is the indication which map to use
// key is the key...
// ch is the channel to register
func RegisterChan(dir int, key string, ch chan *GossipItem) {
  NumberLock[dir].Lock()
  DirectorChans[dir][key] = ch
  NumberLock[dir].Unlock()
}

// SendItem sends the item <it> on the channel indicated by <dir> and <key>
// it needs to obtain a reading lock on the map,
// since maps are not thread safe
func SendItem(dir int, key string, it *GossipItem) (ok bool) {
  NumberLock[dir].RLock()
  ch := DirectorChans[dir][key]
  NumberLock[dir].RUnlock()
  if ch == nil {
    return false
  }
  select {
    // trying the channel, it full, we return false
  case ch <- it:
    ok = true
  default:
    ok = false
  }
  return
}

// DirectItem looks at the message and finds the correct channel to transfer the item to
func DirectItem(item *GossipItem) {
  msg := item.msg
  hd := msg.Header["Via"]
  if len(hd) > 0 {
    str := hd[0]
    m := viaReg.FindStringSubmatch(str)
    if m != nil && len(m) > 1 {
      via := m[1]
      if SendItem(Via, via, item) {
        return
      }
    }
  }
  hd = msg.Header["Call-ID"]
  if len(hd) > 0 {
    str := strings.TrimSpace(hd[0])
    if SendItem(CallId, str, item) {
      return
    }
  }
  m := siptelReg.FindStringSubmatch(msg.SipLine)
  if m != nil && len(m) > 1 {
    telno := m[1]
    if SendItem(Number, telno, item) {
      return
    }
  }
  log.Fatal("don't know where to send ", item)
}

// init gets called during initialization, each init function in each file
func init() {

  NumberLock = make([]sync.RWMutex, DirectoryEnd)
  DirectorChans = make([]director, DirectoryEnd)
  for i := 0; i < DirectoryEnd; i++ {
    DirectorChans[i] = make(director)
  }
  // starting a background goroutine
  go cleanUpDirector()
  // compiling some regular expression
  re, err := regexp.Compile("branch=([^; ]*)")
  if err != nil {
    log.Fatal(err)
  }
  viaReg = re
  re, err = regexp.Compile(":([^@; ]*)")
  if err != nil {
    log.Fatal(err)
  }
  siptelReg = re
}
