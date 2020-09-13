package main

import(
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
)

type GossipTest2 struct {

}

type GossipMsgGrab struct{
  Header string `yaml:"header"`
  Body bool `yaml:"body"`
  RegExp string `yaml:"body"`
  Storage map[string]string `yaml:"store"`
}

type GossipTestMsg struct{
  In string `yaml:"in,omitempty"`
  Out string `yaml:"out,omitempty"`
  To string `yaml:"to,omitempty"`
  Anoa int `yaml:"anoa,omitempty"`
  Bnoa int `yaml:"bnoa,omitempty"`
  Variables map[string]string `yaml:"variables,omitempty"`
  Flags []string `yaml:"flags,omitempty"`
  Headers string `yaml:"headers,omitempty"`
  Rtp int `yaml:"rtp,omitempty"`
  Auto bool `yaml:"auto,omitempty"`
  Delay string `yaml:"delay,omitempty"`
  Test []*GossipTest2 `yaml:"test,omitempty"`
}

type GossipTestCall struct{
  Number string `yaml:"number"`
  Msgs []*GossipTestMsg `yaml:"seq"`
}

type GossipTest struct{
  Name string `yaml:"name"`
  Calls []*GossipTestCall `yaml:"calls"`
}

type TestSuite struct{
  Suite string `yaml:"suite"`
  Tests [] *GossipTest `yaml:"tests"`
}

func GetTestSuite(fn string)(ts *TestSuite){
  ts=new(TestSuite)
  data, err := ioutil.ReadFile(fn)
  if err != nil {log.Fatal(err)}
  err = yaml.Unmarshal(data, ts)
  if err != nil {log.Fatal(err)}
  return
}
