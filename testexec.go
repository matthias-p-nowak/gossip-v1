package main

import (
  "hash/adler32"
  "log"
  "sync"
  "sync/atomic"
)

var (
  TestsUnderExec sync.WaitGroup
  TestProcCount  int32 = 0
  TestTasks            = make(chan *GossipTest, 1)
  executorNumber int32
  TestLocks map[*GossipTest] *sync.Mutex
)

/*
 * A number of parallel working executors are taking tests from main, there are as many executors as there are cfg.Concurrent
 * The number adapts to changes each time a test is done or a new test is commenced.
 */
func TestExecutor() {
  execNumber := atomic.AddInt32(&executorNumber, 1)
  log.Printf("executor %d started\n", execNumber)
  defer log.Printf("executor %d stopped\n", execNumber)
  for {
    // Test if we have to decrease the number of executors
    if atomic.AddInt32(&TestProcCount, 1) > cfg.Concurrent {
      // too many
      log.Printf("executor %d retiring\n", execNumber)
      // take this one back
      atomic.AddInt32(&TestProcCount, -1)
      return
    }
    // running the test
    test := <-TestTasks // fetching one test
    TestLocks[test].Lock()
    run(test)
    TestLocks[test].Unlock()
    // one test done
    TestsUnderExec.Done()
    // this loops ended
    atomic.AddInt32(&TestProcCount, -1)
  }
}

func run(test *GossipTest) {
  var wg sync.WaitGroup
  wg.Add(len(test.Calls))
  for _, call := range test.Calls {
    n := call.Number
    ch := make(chan *GossipItem, 32)
    RegisterChan(Number, n, ch)
    go runCall(call, &wg, ch)
  }
  log.Println("### waiting")
  wg.Wait()
  log.Println("run done")
}

func runCall(call *GossipTestCall, wg *sync.WaitGroup, ch chan *GossipItem) {
  defer wg.Done()
  defer FillChan(ch)
  log.Println("running " + call.Number)
  hadMsg := make(map[uint32]bool)
  // compiling the call
  actions := compileCall(call)
  // doing all the actions
  // TODO, rewrite it using a function for receiving messages with a timeout
  for msg := range ch {
    // if msg is nil check for delayed actions
    if msg != nil {
      if len(actions) == 0 {
        break // nothing more to do
      }
      if msg.msg != nil {
        hash := adler32.Checksum(msg.msg.RawMsg)
        if hadMsg[hash] {
          // just ignoring this message, since it was repeated
          continue
        }
        hadMsg[hash] = true
      }
    }
    var newActions []*Action
    for _, action := range actions {
      a, res := (*action).Examine(msg)
      newActions = append(newActions, a...)
      switch res {
      default:
        log.Printf("don't know what to do with %d", res)
      }
    }
    actions=newActions
  }
}

func compileCall(call *GossipTestCall) (ret []*Action) {
  return
}
