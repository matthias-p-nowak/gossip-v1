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

/*
 * A number of parallel working executors are taking tests from main, there are as many executors as there are cfg.Concurrent
 * The number adapts to changes each time a test is done or a new test is commenced.
 */
func TestRunner() {
	execNumber := atomic.AddInt32(&curRunner, 1)
	if cfg.Verbose >= VerboseTestRunners {
		fmt.Printf("test runner %d started\n", execNumber)
		defer fmt.Printf("test runner %d stopped\n", execNumber)
	}
	for {
		// Test if we have to decrease the number of executors
		if atomic.AddInt32(&TestRunnerCount, 1) > cfg.Concurrent {
			// too many, no log - it is already in the defer list
			// take this one back
			atomic.AddInt32(&TestRunnerCount, -1)
			return
		}
		// running the test
		test := <-TestTasks // fetching one test
		TestLocks[test].Lock()
		tester := Tester{}
		tester.Compile(test)
		tester.Run()
		TestLocks[test].Unlock()
		// one test done
		TestsUnderExec.Done()
		// this loops ended
		atomic.AddInt32(&TestRunnerCount, -1)
	}
}
