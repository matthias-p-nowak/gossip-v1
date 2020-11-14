package main

/*
 * Contains the configuration data
 */
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Config contains the config data provided by YAML
type Config struct {
	Continuous bool     `yaml:"continuous"`
	Loops      int      `yaml:"loops"`
	Rate       int      `yaml:"rate"`
	Concurrent int32    `yaml:"concurrent"`
	Local      []string `yaml:"local"`
	Remote     string   `yaml:"remote"`
	Verbose    int      `yaml:"verbose"`
}

// GetConfig reads the content of <fn> and returns configuration
func GetConfig(fn string) (cfg *Config, err error) {
	cfg = new(Config)
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println("\nERROR not a valid config file, use something like the following")
		io.Copy(os.Stdout, GetStored("snippets/gossip.cfg"))
		return
	}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatal(err)
	}
	/*
	   if(*verbose >7 ){
	   data, err=yaml.Marshal(cfg)
	   fmt.Println(string(data))
	   }
	*/
	return
}
