package main

import(
  "hash/adler32"
  "sync"
)

type Tester struct {
  hadMsg map[uint32]bool
  wg sync.WaitGroup
  actions[][]*TestAction
}

func (t *Tester) Compile(test *GossipTest){

}

func (t* Tester) Run(){
  l:=len(t.actions)
  t.wg.Add(l)
  for r:=0;r<l;r++ {
    go t.Runner(r)
  }
  t.wg.Wait()
}

// SIP can send the same message several times
func (t* Tester) CheckNew(msg *GossipItem) bool {
  if msg != nil {
   if msg.msg != nil {
      hash := adler32.Checksum(msg.msg.RawMsg)
      if t.hadMsg[hash] {
        return false
      } else {
        t.hadMsg[hash] = true
        return true
      }
    }
  }
  return false
}

func (t* Tester) Runner(r int){
  // do the work
  t.wg.Done()
}

type TestAction interface{
  Examine(t *Tester) (ta *TestAction)
}
