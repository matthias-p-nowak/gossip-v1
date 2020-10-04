package main

import (
  "log"
  "time"
)

const (
  ActionSkip = iota
  ActionSuccess
  ActionFailed
)

type Action interface {
  // setup from data
  Compile(tp *TestParty, m *GossipTestMsg)
  // do a normal execution including reading from channel
  Run() (next Action, result int)
  // is this a message this can handle?
  IsMatch(gi *GossipItem) bool
  // got an item and do the relevant stuff
  Execute(gi *GossipItem) (next Action, result int)
  // add a single action to the next
  SetNext(next Action) (ok bool)
}

type DefaultAction struct {
  nextAction Action
  tp         *TestParty
  msg        *GossipTestMsg
}

func (da *DefaultAction) SetNext(next Action) bool {
  if da.nextAction != nil {
    return false
  }
  da.nextAction = next
  return true
}

func (da *DefaultAction) IsMatch(gi *GossipItem) bool {
  return false
}

func (da *DefaultAction) Execute(gi *GossipItem) (next Action, result int) {
  log.Fatal("this should never be called - something wrong")
  return da.nextAction, ActionFailed
}
func (da *DefaultAction) Compile(tp *TestParty, msg *GossipTestMsg) {
  da.tp=tp
  da.msg=msg
}

type DelayAction struct {
  DefaultAction
  Duration   time.Duration
}

func (da *DelayAction) Compile(tp *TestParty, msg *GossipTestMsg) {
  da.tp = tp
  dur, err := time.ParseDuration(msg.Delay)
  if err != nil {
    log.Fatal(err, "\n", tp.Call.String())
  }
  da.Duration = dur
}

func (da *DelayAction) Run() (next Action, result int) {
  ch := time.After(da.Duration)
  if cfg.Verbose > 9 {
    log.Printf("waiting for %s for %s\n", da.Duration.String(), da.tp.Call.Number)
  }
  for {
    select {
    case <-ch:
      if cfg.Verbose > 9 {
        log.Printf("  waiting ended for %s\n", da.tp.Call.Number)
      }
      return da.nextAction, ActionSuccess
    case gi := <-da.tp.ch:
      if gi != nil {
        if da.tp.tester.CheckNew(gi) {
          return da.tp.CheckOptional(gi)
        }
      }
    }
  }
}


type SendInvite struct {
  DefaultAction
}

func (si *SendInvite) Run() (next Action, result int) {
  return
}


