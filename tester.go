package main

import (
  "log"
  "strings"
  "sync"
)

// compiled for the whole test
type Tester struct {
  hadMsg  map[uint32]bool
  wg      sync.WaitGroup
  parties []*TestParty
}

// for one party
type TestParty struct {
  tester          *Tester
  ch              chan *GossipItem
  Call            *GossipTestCallParty
  actions         []Action
  optionalActions []Action
}

func (t *Tester) Insert(p *TestParty, msg *GossipTestMsg, ac Action) {
  l := len(p.actions)
  if l > 0 {
    ac.SetNext(p.actions[l-1])
  }
  ac.Compile(p, msg)
  p.actions = append(p.actions, ac)
}

func (t *Tester) Compile(test *GossipTest) {
  for i, c := range test.Calls {
    p := new(TestParty)
    p.tester = t
    t.parties = append(t.parties, p)
    p.ch = make(chan *GossipItem, 8)
    RegisterChan(Number, c.Number, p.ch)
    p.Call = c
    for j, msg := range c.Msgs {
      if cfg.Verbose > 5 {
        log.Printf("compiling %s.%d.%d: %s\n", test.Name, i, j, c.Number)
      }
      if len(msg.Delay) > 0 {
        if cfg.Verbose > 7 {
          log.Printf("  adding delay action for %s\n", msg.Delay)
        }
        da := new(DelayAction)
        t.Insert(p, msg, da)
      }
      if len(msg.In) > 0 && len(msg.Out) > 0 {
        log.Fatal("don't use In and Out in the same message:", c.String())
      }
      switch {
      case len(msg.Out) > 0:
        switch strings.ToUpper(msg.Out) {
        case "INVITE":
          t.Insert(p, msg, new(SendInvite))
        default:
          log.Fatal("outgoing request " + msg.Out + " is unknown")
        }
      default:
        log.Fatal("don't know what to do with this:\n" + msg.String() + "\n")
      }
    }
  }
}

func (t *Tester) Run() {
  l := len(t.parties)
  t.wg.Add(l)
  for r := 0; r < l; r++ {
    go t.parties[r].Runner()
  }
  t.wg.Wait()
}

// SIP can send the same message several times
func (t *Tester) CheckNew(gi *GossipItem) bool {
  if gi != nil {
      if t.hadMsg[gi.Hash] {
        return false
      } else {
        t.hadMsg[gi.Hash] = true
        return true
      }
  }
  return false
}

func (tp *TestParty) Runner() {
  // TODO: do the work
  if len(tp.actions) > 0 {
    ac := tp.actions[0]
    for ac != nil {
      var res int
      ac, res = ac.Run()
      switch res {
      default:
        log.Fatal("no idea how to handle result ", res)
      }
    }
  } else {
    log.Fatal("no action for ", tp.Call.String())
  }
  tp.tester.wg.Done()
}

func (tp *TestParty) CheckOptional(gi *GossipItem) (next Action, result int) {
  // TODO: check optional message items
  log.Fatal("implement function")
  return
}
