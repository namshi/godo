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
	// The IP address of the server
	Address string
	// The user that will login
	// into the server
	User string
	// An optional host that will
	// be used as SSH tunnel to the
	// current server
	Tunnel string
}

// Represents a command
// to run on a server or
// group of servers.
type Command struct {
	// The target on which we are going
	// to execute the current command: can
	// be a server or a group
	Target string
	// A string representing the command
	// to execute
	Exec string
	// A description of the command,
	// used in the CLI
	Description string
}

// Represents the whole configuration
// file.
type Config struct {
	// The file to use to verify
	// known hosts
	Hostfile string
	// Timeout of the SSH connections
	Timeout int
	// List of servers
	Servers map[string]Server
	// List of commands
	Commands map[string]Command
	// List of groups
	Groups map[string][]string
	// The raw string containing the
	// whole configuration, in YAML
	// format
	Raw string
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

	if c.Timeout == 0 {
		c.Timeout = 5
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
