// Package used to be able
// to remotely execute commands
// via SSH.
//
// Most of this stuff relies on
// CoreOS's SSH package (https://github.com/coreos/fleet/tree/master/ssh).
package ssh

import (
	"time"

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

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Creates a new SSH client based on the
// given configuration.
func CreateClient(config *Config) *ssh.SSHForwardingClient {
	hostfile := ssh.NewHostKeyFile(config.Hostfile)
	checker := ssh.NewHostKeyChecker(hostfile)

	if config.Tunnel != "" {
		client, err := ssh.NewTunnelledSSHClient(config.User, config.Tunnel, config.Address, checker, true, config.Timeout)
		handleError(err)

		return client
	}

	client, err := ssh.NewSSHClient(config.User, config.Address, checker, true, config.Timeout)
	handleError(err)

	return client
}
