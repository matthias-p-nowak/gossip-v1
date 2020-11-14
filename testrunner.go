package main

import (
  "fmt"
  "sync"
  "sync/atomic"
)

var (
  TestsUnderExec  sync.WaitGroup
  TestRunnerCount int32 = 0
  TestTasks             = make(chan *GossipTest, 1)
  curRunner       int32
  TestLocks       = make(map[*GossipTest]*sync.Mutex)
)

// TestRunner is one of many test executors
// A number of parallel working executors are taking tests from main, there are as many executors as there are cfg.Concurrent
// The number adapts to changes each time a test is done or a new test is commenced.
func TestRunner() {
  // giving this one a unique number
  execNumber := atomic.AddInt32(&curRunner, 1)
  if cfg.Verbose >= VerboseTestRunners {
    fmt.Printf("test runner %d started\n", execNumber)
    defer fmt.Printf("test runner %d stopped\n", execNumber)
  }
  // main loop
  for {
    // Test if we have to decrease the number of executors
    if atomic.AddInt32(&TestRunnerCount, 1) > cfg.Concurrent {
      // too many, no log - it is already in the defer list
      // take this one back
      atomic.AddInt32(&TestRunnerCount, -1)
      // and end this one
      return
    }
    // running the test
    test := <-TestTasks // fetching one test
    // don't run this test several times simultaneously
    TestLocks[test].Lock()
    tester := Tester{}
    tester.CompileTest(test)
    tester.RunTest()
    TestLocks[test].Unlock()
    // one test done - it was added by the loop that feeds the TestTasks channel
    TestsUnderExec.Done()
    // this loops ended
    atomic.AddInt32(&TestRunnerCount, -1)
  }
}
