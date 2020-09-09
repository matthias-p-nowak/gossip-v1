package main

import(
  "sync"
  "time"
)

const (
  Number = iota
  CallId
  Via
  DirectoryEnd
)

type director map[string]chan *GossipItem

var (
  DirectorChans []director
  NumberLock []sync.RWMutex
)



func cleanUpDirector(){
  for{
    for i:=0;i<DirectoryEnd;i++{
      time.Sleep(time.Second)
      k2d:=[]string{}
      NumberLock[i].RLock()
        for k,v:=range DirectorChans[i]{
          select {
            case v <- nil:
              // everything ok
            default:
              // channel is filled == no reader left
              k2d=append(k2d,k)
          }
        }
      NumberLock[i].RUnlock()
      NumberLock[i].Lock()
        for _,k:=range k2d {
          delete(DirectorChans[i],k)
        }
      NumberLock[i].Unlock()
      }
  }
}

func RegisterChan(dir int, key string, ch chan *GossipItem){
  NumberLock[dir].Lock()
  DirectorChans[dir][key]=ch
  NumberLock[dir].Unlock()
}


func SendItem(dir int, key string, it *GossipItem) (ok bool) {
  NumberLock[dir].RLock()
  ch:=DirectorChans[dir][key]
  NumberLock[dir].RUnlock()
  if ch == nil {
    return false
  }
  select {
  case ch <- it:
      ok=true
  default:
    ok=false
  }
  return
}

func init(){
  NumberLock=make([]sync.RWMutex,DirectoryEnd)
  DirectorChans=make([]director,DirectoryEnd)
  for i:=0;i<DirectoryEnd;i++{
    DirectorChans[i]=make(director)
  }
  go cleanUpDirector()
}
