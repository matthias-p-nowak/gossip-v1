package main

/*
 * The data structures that represent test input data
 */
import (
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "log"
)

/*
 * TODO: fill with tests that are run on the received messages
 */
type GossipTest2 struct {
}

/*
 * for grabbing parts of messages for further use
 */
type GossipMsgGrab struct {
  Header  string            `yaml:"header"`
  Body    bool              `yaml:"body"`
  RegExp  string            `yaml:"body"`
  Storage map[string]string `yaml:"store"`
}

/*
 * Entry for a single message
 */
type GossipTestMsg struct {
  // Name of this action for back reference
  Alias string `yaml:"alias,omitempty"`
  // single delay in that scenario
  Delay string `yaml:"delay,omitempty"`
  // Message related
  // in/out: Invite...10,200
  In  string `yaml:"in,omitempty"`
  Out string `yaml:"out,omitempty"`
  // for outgoing invites
  To   string `yaml:"to,omitempty"`
  Anoa int    `yaml:"anoa,omitempty"`
  Bnoa int    `yaml:"bnoa,omitempty"`
  // prefills the header structure
  Headers string `yaml:"headers,omitempty"`
  // flags like 100rel
  Flags []string `yaml:"flags,omitempty"`
  // additional definition of variables used in templates
  Variables map[string]string `yaml:"variables,omitempty"`
  // number of rtp packets before succeeding to next item
  Rtp int `yaml:"rtp,omitempty"`
  // fullfil the dialog and end it according to dialog state
  Continue bool `yaml:"continue,omitempty"`
  // list of tests to be carried out
  Test []*GossipTest2 `yaml:"test,omitempty"`
  // grabbing parts of the message and storing in variables
  Grabs []*GossipMsgGrab `yaml:"grab,omitempty"`
  // an optional message that changes the course
  Optional bool `yaml:"optional,omitempty"`
}

// on call party
type GossipTestCallParty struct {
  test_bl *GossipTest
  // number used in FROM for invites, or TO for getting calls to this party
  // must be exactly as in the receiving invite
  Number string           `yaml:"number"`
  Msgs   []*GossipTestMsg `yaml:"seq"`
}

// One single test
// TODO: add number of runs
// TODO: add global tests
type GossipTest struct {
  suite_bl *TestSuite
  Name     string                 `yaml:"name"`
  Calls    []*GossipTestCallParty `yaml:"calls"`
}

// a set of tests, used to group tests
type TestSuite struct {
  Suite string        `yaml:"suite"`
  Tests []*GossipTest `yaml:"tests"`
}

// Reading the test file and add a suite
func GetTestSuite(fn string) (ts *TestSuite) {
  ts = new(TestSuite)
  data, err := ioutil.ReadFile(fn)
  if err != nil {
    log.Fatal(err)
  }
  err = yaml.Unmarshal(data, ts)
  if err != nil {
    log.Fatal(err)
  }
  for _, t := range ts.Tests {
    t.suite_bl = ts
    for _, cp := range t.Calls {
      cp.test_bl = t
    }
  }
  return
}

func (cp *GossipTestCallParty) String() string {
  t := cp.test_bl
  ts := t.suite_bl
  s1 := fmt.Sprintf("# %s -> %s -> %s \n", ts.Suite, t.Name, cp.Number)
  data, err := yaml.Marshal(cp)
  if err != nil {
    log.Fatal(err)
  }
  return s1 + string(data)
}

func (gtm *GossipTestMsg) String() string {
  data, err := yaml.Marshal(gtm)
  if err != nil {
    log.Fatal(err)
  }
  return string(data)
}
