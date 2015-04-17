// This package is responsible for parsing
// a configuration file and generate
// Go structures.
package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

// Represents a server on which
// we can connect to.
type Server struct {
	Address string
	User    string
	Tunnel  string
}

// Represents a command
// to run on a server or
// group of servers.
type Command struct {
	Target      string
	Exec        string
	Description string
}

// Represents the whole configuration
// file.
type Config struct {
	Hostfile string
	Servers  map[string]Server
	Commands map[string]Command
	Raw      string
	Groups   map[string][]string
}

// Parses a configuration file and returns
// a Go structure that represents it.
func Parse(file string) Config {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		panic(err)
	}

	c := Config{}

	err = yaml.Unmarshal([]byte(content), &c)

	if err != nil {
		panic(err)
	}

	c.Raw = string(content)

	return c
}
