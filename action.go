package main

import (
	"log"
	"time"
)

const (
	// Results of Actions
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
	Run() (next int, result int)
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
	msg *GossipTestMsg
	// position this action is stored in the tp.actions slice
	pos int
}

// IsMatch per default false, must be overriden by optional actions
func (da *DefaultAction) IsMatch(gi *GossipItem) bool {
	return false
}

// DefaultNext - next action is usually the next one to execute
func (da *DefaultAction) DefaultNext() int {
	return da.pos + 1
}

// Execute deals with an incoming item, called within Run
func (da *DefaultAction) Execute(gi *GossipItem) (next int, result int) {
	log.Fatal("this should never be called - something wrong")
	next = da.DefaultNext()
	result = ActionFailed
	return
}

// Compile prepares the execution of this test,
func (da *DefaultAction) Compile(tp *TestParty, msg *GossipTestMsg) {
	da.msg = msg
}

// GetTransaction returns the transaction for this action, if there is one
func (da *DefaultAction) GetTransaction() (tr *GossipTransaction) {
	return
}

// DelayAction for simple delays
type DelayAction struct {
	DefaultAction
	Duration time.Duration
}

// Compile extracts the delay
func (da *DelayAction) Compile(tp *TestParty, msg *GossipTestMsg) {
	da.tp = tp
	dur, err := time.ParseDuration(msg.Delay)
	if err != nil {
		log.Fatal(err, "\n", tp.CallParty.String())
	}
	da.Duration = dur
}

// Run executes the delay
func (da *DelayAction) Run() (next int, result int) {
	ch := time.After(da.Duration)
	if cfg.Verbose > 9 {
		log.Printf("waiting for %s for %s\n", da.Duration.String(), da.tp.CallParty.Number)
	}
	// need endless for-loop due to nil-items
	for da.tp.tester.running {
		select {
		// timed out
		case <-ch:
			if cfg.Verbose > 9 {
				log.Printf("  waiting ended for %s\n", da.tp.CallParty.Number)
			}
			return da.DefaultNext(), ActionSuccess
			// got something from the call party related channel
		case gi := <-da.tp.ch:
			if gi != nil {
				// ignore the nil messages
				if da.tp.tester.CheckNew(gi) {
					return da.tp.CheckOptional(gi)
				}
			}
		}
	}
	next = -1
	result = ActionFailed
	return
}

// SendInvite action sends invites
type SendInvite struct {
	DefaultAction
	tr *GossipTransaction
}

// Compile prepares the SendInvite action
func (si *SendInvite) Compile(tp *TestParty, m *GossipTestMsg) {
	if len(tp.actions) == 0 {
		// the first invite
		c := new(GossipCall)
		d := c.NewDialog()
		t := d.NewTransaction()
		si.tr = t
	}
}

// GetTransaction returns the created transaction
func (si *SendInvite) GetTransaction() (tr *GossipTransaction) {
	tr = si.tr
	return
}

// Run executes the SendInvite action
func (si *SendInvite) Run() (next int, result int) {
	next = -1
	result = ActionSuccess
	return
}
