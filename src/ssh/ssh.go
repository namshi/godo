// Package used to be able
// to remotely execute commands
// via SSH.
//
// Most of this stuff relies on
// CoreOS's SSH package (https://github.com/coreos/fleet/tree/master/ssh).
package ssh

import (
	"time"

	gossh "github.com/coreos/fleet/Godeps/_workspace/src/golang.org/x/crypto/ssh"
	ssh "github.com/coreos/fleet/ssh"
)

// A simple configuration to connect
// to an SSH host.
// If a Tunnel is provided, it will be
// used as bastion host to connect to
// the "real" address.
type Config struct {
	Hostfile string
	Address  string
	Tunnel   string
	User     string
	Alias    string
	Timeout  time.Duration
}

// Creates a new SSH session to run
// commands with the given config
// parameters.
func CreateSession(config *Config) (*gossh.Session, error) {
	hostfile := ssh.NewHostKeyFile(config.Hostfile)
	checker := ssh.NewHostKeyChecker(hostfile)

	if config.Tunnel != "" {
		client, _ := ssh.NewTunnelledSSHClient(config.User, config.Tunnel, config.Address, checker, true, config.Timeout)

		return client.NewSession()
	}

	// TODO: fix this
	return nil, nil
}
