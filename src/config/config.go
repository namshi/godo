// This package is responsible for parsing
// a configuration file and generate
// Go structures.
package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
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

// Tries to read the contents of the file
// trying to locate it in different  directories.
// ie. ~/godo.yml | ./godo.yml
func getFileContents(file string) []byte {
	user, _ := user.Current()
	wd, _ := os.Getwd()
	homePath := path.Join(user.HomeDir, file)
	wdPath := path.Join(wd, file)
	paths := []string{wdPath, homePath, file}

	for _, path := range paths {
		content, err := ioutil.ReadFile(path)

		if err == nil {
			return content
		}
	}

	log.Fatalf("Unable to read config file (tried %s)", strings.Join(paths, ", "))
	return []byte{}
}

// Validates the parsed configuration,
// providing default values if anything
// is missing.
func validate(c *Config) {
	if c.Hostfile == "" {
		c.Hostfile = path.Join(os.TempDir(), "known_hosts_godo")
	}
}

// Parses a configuration file and returns
// a Go structure that represents it.
func Parse(file string) Config {
	content := getFileContents(file)
	c := Config{}

	err := yaml.Unmarshal([]byte(content), &c)

	if err != nil {
		panic(err)
	}

	c.Raw = string(content)
	validate(&c)

	return c
}
