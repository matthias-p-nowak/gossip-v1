package main

import(
  "log"
  "sync"
  "time"
)

var (
  director sync.Map
  cleaned int
)


func checkAndDelete(key interface{}, val interface{}) bool {
  keyStr:=key.(string)
  valChan:=val.(chan *GossipItem)
  select{
    case valChan <- nil:
    default:
      cleaned++
      director.Delete(keyStr)
  }
  return true
}

func cleanUpDirector(){
  for{
    time.Sleep(time.Second)
    cleaned=0
    director.Range(checkAndDelete)
    if cleaned >0 {
      log.Printf("cleaned %d\n",cleaned)
    }
  }
}

func RegisterChan(key string, ch chan *GossipItem){
  director.Store(key,ch)
}

func SendItem(key string, it *GossipItem) (ok bool) {
  if val,ok:=director.Load(key);ok {
    ch:=val.(chan *GossipItem)
    ch <- it
  }
  return
}

func init(){
  log.Println("starting cleanup")
  go cleanUpDirector()
}
