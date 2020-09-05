package main

import (
  "flag"
  "log"
  "os"
  "os/signal"
  "path/filepath"
)

//go:generate go run scripts/go-bin.go -o snippets.go snippets

var (
  config *Config
  gossipRunning = true
  testSuites []*TestSuite
  verbose *int
)

func setup() {

}

func handleSignals() {
  signals := make(chan os.Signal, 1)
  signal.Notify(signals)
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
  testSuites=append(testSuites,ts)
  return nil
}

func main() {
  defer log.Println("gossip is done")
  log.SetFlags(log.LstdFlags | log.Lshortfile)
  log.Println("gossip started")
  cfgFile := flag.String("c", "gossip.cfg", "the configuration for gossip")
  verbose= flag.Int("v",0,"verbosity of gossip for testing; higher means more")
  flag.Parse()
  config,err := GetConfig(*cfgFile)
  if err != nil { log.Fatal(err)}
  setup()
  for _, arg:=range flag.Args() {
    log.Println("investigating "+arg)
    filepath.Walk(arg,parseTests)
  }
  go handleSignals()
  if config.Loops < 1 {
    config.Loops = 1
  }
  for i := 0; i < config.Loops; i++ {
    log.Printf("loop %d\n",i)
  }
}
