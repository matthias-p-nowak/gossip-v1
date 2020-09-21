package main

import (
	"runtime"
	"sync"
	"time"
)

// used as index for the different maps
const (
	Number = iota
	CallId
	Via
	DirectoryEnd
)

// maps a certain tag to a channel for gossip items
type director map[string]chan *GossipItem

var (
	// the maps from tags to strings
	DirectorChans []director
	// the related mutex for access serialization
	NumberLock []sync.RWMutex
)

// when the goroutines serving the channels end, they might have been registered several times
func cleanUpDirector() {
	for {
		for i := 0; i < DirectoryEnd; i++ {
			time.Sleep(time.Second)
			k2d := []string{}
			NumberLock[i].RLock()
			for k, v := range DirectorChans[i] {
				select {
				case v <- nil:
					// everything ok
				default:
					// channel is filled == no reader left
					k2d = append(k2d, k)
				}
			}
			NumberLock[i].RUnlock()
			NumberLock[i].Lock()
			for _, k := range k2d {
				delete(DirectorChans[i], k)
			}
			NumberLock[i].Unlock()
			runtime.Gosched()
		}
	}
}

func RegisterChan(dir int, key string, ch chan *GossipItem) {
	NumberLock[dir].Lock()
	DirectorChans[dir][key] = ch
	NumberLock[dir].Unlock()
}

func FillChan(ch chan *GossipItem) {
	for {
		select {
		case ch <- nil:
			// ok
		default:
			return
		}
	}
}

func SendItem(dir int, key string, it *GossipItem) (ok bool) {
	NumberLock[dir].RLock()
	ch := DirectorChans[dir][key]
	NumberLock[dir].RUnlock()
	if ch == nil {
		return false
	}
	select {
	case ch <- it:
		ok = true
	default:
		ok = false
	}
	return
}

func init() {
	NumberLock = make([]sync.RWMutex, DirectoryEnd)
	DirectorChans = make([]director, DirectoryEnd)
	for i := 0; i < DirectoryEnd; i++ {
		DirectorChans[i] = make(director)
	}
	go cleanUpDirector()
}
