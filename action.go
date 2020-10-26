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

// Action is the basic interface for all actions
type Action interface {
  // setup from data
  // TODO - remove compile, make special functions for each one
  Compile(tp *TestParty, m *GossipTestMsg)
  // do a normal execution including reading from channel, when reading items from channel, call Execute
  Run(tp *TestParty) (next int, result int)
  // is this a message this can handle?
  IsMatch(gi *GossipItem) bool
  // got an item and do the relevant stuff
  Execute(gi *GossipItem) (next int, result int)
  GetTransaction() (tr *GossipTransaction)
}

// DefaultAction combines the common parts for all actions
type DefaultAction struct {
  // backlink
  tp *TestParty
  // data for this action from the test suite
  msg        *GossipTestMsg
  // position this action is stored in the tp.actions slice
  pos int
}

// IsMatch per default false, must be overriden by optional actions
func (da *DefaultAction) IsMatch(gi *GossipItem) bool {
  return false
}


func (da *DefaultAction) DefaultNext(tp *TestParty) int {
  return da.pos+1
}
// 
func (da *DefaultAction) Execute(gi *GossipItem) (next int, result int) {
  log.Fatal("this should never be called - something wrong")
  next=da.DefaultNext()
  result=ActionFailed
  return 
}
func (da *DefaultAction) Compile(tp *TestParty, msg *GossipTestMsg) {
  da.msg = msg
}

func (da *DefaultAction) GetTransaction() (tr *GossipTransaction){
  return
}

type DelayAction struct {
  DefaultAction
  Duration time.Duration
}

func (da *DelayAction) Compile(tp *TestParty, msg *GossipTestMsg) {
  da.tp = tp
  dur, err := time.ParseDuration(msg.Delay)
  if err != nil {
    log.Fatal(err, "\n", tp.CallParty.String())
  }
  da.Duration = dur
}

func (da *DelayAction) Run() (next Action, result int) {
  ch := time.After(da.Duration)
  if cfg.Verbose > 9 {
    log.Printf("waiting for %s for %s\n", da.Duration.String(), da.tp.CallParty.Number)
  }
  for {
    select {
    case <-ch:
      if cfg.Verbose > 9 {
        log.Printf("  waiting ended for %s\n", da.tp.CallParty.Number)
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
  tr *GossipTransaction
}

func (si *SendInvite) Compile(tp *TestParty, m *GossipTestMsg){
  if len(tp.actions) == 0 {
    // the first invite
    c:=new(GossipCall)
    d:=c.NewDialog()
    t:=d.NewTransaction()
    si.tr=t
  }
}

func (si *SendInvite)GetTransaction()(tr *GossipTransaction){
  tr=si.tr
  return
}

func (si *SendInvite) Run() (next Action, result int) {
  next=si.nextAction
  return
}
