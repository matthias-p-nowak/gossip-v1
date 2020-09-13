package main

const (
  ActionSkip = iota
  ActionSuccess
  ActionFailed
)

type Action interface {
  Activate()
  Examine(msg *GossipItem) (next []*Action, result int)
}
