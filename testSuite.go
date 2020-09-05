package main

import(
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
)

type GossipTestService struct{
  Number string `yaml:"num,omitempty"`
  Place string `yaml:"place,omitempty"`
}
type GossipTestMsg struct{
  In string `yaml:"in,omitempty"`
  Out string `yaml:"out,omitempty"`
  To string `yaml:"to,omitempty"`
  Service GossipTestService `yaml:"service,omitempty"`
  Anoa int `yaml:"anoa,omitempty"`
  Bnoa int `yaml:"bnoa,omitempty"`
  Flags []string `yaml:"flags,omitempty"`
  Headers string `yaml:"headers,omitempty"`
  Rtp int `yaml:"rtp,omitempty"`
  Auto bool `yaml:"auto,omitempty"`
  Delay int `yaml:"delay,omitempty"`
}

type GossipTestCall struct{
  Number string `yaml:"number"`
  Msgs []GossipTestMsg `yaml:"seq"`
}

type GossipTestRun struct{
  Name string `yaml:"name"`
  Calls []GossipTestCall `yaml:"calls"`
}

type GossipTest struct {
  Name string `yaml:"name"`
  Runs []GossipTestRun `yaml:"runs"`
}

type TestSuite struct{
  Suite string `yaml:"suite"`
  Tests []GossipTest `yaml:"tests"`
}

func GetTestSuite(fn string)(ts *TestSuite){
  ts=new(TestSuite)
    data, err := ioutil.ReadFile(fn)
  if err != nil {
    log.Fatal(err)
    }
  err = yaml.Unmarshal(data, ts)
  if err != nil {log.Fatal(err)}
  if(*verbose >7 ){
  data, err=yaml.Marshal(ts)
  fmt.Println(string(data))
}
  return
}
