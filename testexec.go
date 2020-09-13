package main

import (
  "log"
  "sync"
  "sync/atomic"
)

var (
  TestsUnderExec sync.WaitGroup
  TestProcCount int32=0
  TestTasks = make(chan *GossipTest, 1)
  executorNumber int32
)

func TestExecutor(){
  execNumber:=atomic.AddInt32(&executorNumber,1)
  log.Printf("executor %d started\n",execNumber)
  defer log.Printf("executor %d stopped\n",execNumber)
  for {
    // are we superfluous?
       if atomic.AddInt32(&TestProcCount,1) > cfg.Concurrent {
      log.Printf("executor %d retiring\n",execNumber)
      atomic.AddInt32(&TestProcCount,-1)
      return
    }
    // running the test
    test:= <- TestTasks
    test.Lock.Lock()
    run(test)
    test.Lock.Unlock()
    // one test done
    TestsUnderExec.Done()
    // this loops ended
    atomic.AddInt32(&TestProcCount,-1)
  }
}

func run(test *GossipTest){
  var wg sync.WaitGroup
  wg.Add(len(test.Calls))
  for _,call := range test.Calls{
    n:=call.Number
    ch:=make(chan *GossipItem,32)
    RegisterChan(Number, n, ch)
    go runCall(call,&wg,ch)
  }
  log.Println("### waiting")
  wg.Wait()
  log.Println("run done")
}

func runCall(call *GossipTestCall, wg *sync.WaitGroup, ch chan *GossipItem){
  defer wg.Done()
  defer FillChan(ch)
  // compiling the call
  log.Println("running "+call.Number)
}
