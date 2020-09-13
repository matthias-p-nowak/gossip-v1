package main

import (
  "flag"
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
  "os"
  "os/signal"
  "path/filepath"
  "syscall"
  "sync"
  "sync/atomic"
)

//go:generate go run scripts/go-bin.go -o snippets.go snippets

var (
  cfg *Config
  gossipRunning = true
  testSuites []*TestSuite
)

func setup() {
  LimiterInit()
  for _,provider:=range cfg.Local {
    NewProvider(provider)
  }
}

func handleSignals() {
  signals := make(chan os.Signal, 1)
  signal.Notify(signals,syscall.SIGHUP)
  for s := range signals {
    log.Println("Got signal:", s)
  }
}

func parseTests(fn string, info os.FileInfo, err error) error{
  if ! info.Mode().IsRegular() {
    log.Println("skipping "+fn)
    return nil
  }
  log.Println("parsing "+fn)
  ts:=GetTestSuite(fn)
  if(cfg.Verbose >7 ){
  data, err:=yaml.Marshal(ts)
  if err != nil {log.Fatal(err)}
  fmt.Println(string(data))
}
  testSuites=append(testSuites,ts)
  return nil
}

func main() {
  defer log.Println("gossip is done")
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  log.Println("gossip started")
  cfgFile := flag.String("c", "gossip.cfg", "the configuration for gossip")
  verbose:= flag.Int("v",-1,"verbosity of gossip for testing; higher means more")
  flag.Parse()
  var err error
  cfg,err = GetConfig(*cfgFile)
  if err != nil { log.Fatal(err)}
  if *verbose>=0 {
    cfg.Verbose=*verbose
  }
  setup()
  for _, arg:=range flag.Args() {
    log.Println("investigating "+arg)
    filepath.Walk(arg,parseTests)
  }
  go handleSignals()
  if cfg.Loops < 1 {
    cfg.Loops = 1
  }
  // creating Mutex for each test
  for _,ts:=range testSuites{
    for _,t:=range ts.Tests {
      TestLocks[t]=new(sync.Mutex)
    }
  }
  for i := 0; i < cfg.Loops; i++ {
    if cfg.Continuous {
      i=0
    }
    log.Printf("loop %d\n",i)
    for tsi,ts:=range testSuites{
      log.Printf("test suite(%d): %s\n",tsi,ts.Suite)
      for ti,test:=range ts.Tests {
        log.Printf(" test(%d): %s\n",ti,test.Name)
        TestsUnderExec.Add(1)
        if atomic.LoadInt32(&TestProcCount) < cfg.Concurrent {
          log.Println("starting a new test executor")
          go TestExecutor()
        }
        TestTasks <- test
      }
    }
  }
  log.Println("waiting for tests to finish")
  TestsUnderExec.Wait()
}
