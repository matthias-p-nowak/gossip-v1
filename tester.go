package main

import (
	"log"
	"regexp"
	"strings"
	"sync"
)

// Tester - main structure for executing a certain test
// compiled for the whole test
type Tester struct {
	// a map of hash values of messages received before
	hadMsg map[uint32]bool
	// wg waits on all call parties to be finished
	wg sync.WaitGroup
	// parties == all end points involved in a test call
	parties []*TestParty
	// the initiating call party should set this to false
	running bool
}

// TestParty - for one party
type TestParty struct {
	// backlink
	tester *Tester
	// dedicated channel for this call party
	ch chan *GossipItem
	// input data
	CallParty *GossipTestCallParty
	// the normal actions that the RunTest goes through
	actions []Action
	// if the incoming message does not match, these are tried
	optionalActions []Action
	// map of alias actions
	actionMap map[string]int
}

var (
	// respCodeReg is the regular expression for 3 digits
	respCodeReg *regexp.Regexp
)

func init() {
	// initializing respCodeReg
	r, err := regexp.Compile("([0-9]{3})")
	if err != nil {
		log.Fatal(err)
	}
	respCodeReg = r
}

// Insert takes ac and does some common actions
func (t *Tester) Insert(p *TestParty, msg *GossipTestMsg, ac Action) {
	// default
	l := len(p.actions)
	p.actions = append(p.actions, ac)
	if len(msg.Alias) > 0 {
		p.actionMap[msg.Alias] = l
	}
}

// creates a new TestParty and initializes it
func (t *Tester) createTestParty(c *GossipTestCallParty) (tp *TestParty) {
	tp = new(TestParty)
	tp.tester = t
	tp.actionMap = make(map[string]int)
	t.parties = append(t.parties, tp)
	tp.ch = make(chan *GossipItem, 8)
	tp.CallParty = c
	return
}

// CompileTest is the first phase, where datastructures are created and connections established
func (t *Tester) CompileTest(test *GossipTest) {
	//
	for i, c := range test.CallParties {
		p := t.createTestParty(c)
		RegisterChan(Number, c.Number, p.ch)

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
				m := respCodeReg.FindStringSubmatch(msg.Out)
				if m != nil {

				} else {
					switch strings.ToUpper(msg.Out) {
					case "INVITE":
						t.Insert(p, msg, new(SendInvite))
					default:
						log.Fatal("outgoing request " + msg.Out + " is unknown")
					}
				}
			default:
				log.Fatal("don't know what to do with this:\n" + msg.String() + "\n")
			}
		}
	}
}

// RunTest
func (t *Tester) RunTest() {
	t.running = true
	l := len(t.parties)
	t.wg.Add(l)
	for r := 0; r < l; r++ {
		go t.parties[r].RunTest()
	}
	t.wg.Wait()
}

// CheckNew checks if this message was received before, this one being a retransmission
// SIP can send the same message several times in retransmissions
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

func (tp *TestParty) RunTest() {
	defer tp.tester.wg.Done()
	// TODO: do the work
	if len(tp.actions) > 0 {
		pos := 0
		for pos >= 0 && pos <= len(tp.actions) {
			var res int
			if !tp.tester.running {
				break
			}
			pos, res = tp.actions[pos].Run()
			switch res {
			default:
				log.Fatal("no idea how to handle result ", res)
			}
		}
	} else {
		log.Fatal("no action for ", tp.CallParty.String())
	}
}

func (tp *TestParty) CheckOptional(gi *GossipItem) (next int, result int) {
	// TODO: check optional message items
	next = -1
	result = ActionSuccess
	log.Fatal("implement function")
	return
}
