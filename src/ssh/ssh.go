// Package used to be able
// to remotely execute commands
// via SSH.
//
// Most of this stuff relies on
// CoreOS's SSH package (https://github.com/coreos/fleet/tree/master/ssh).
package ssh

import (
	gossh "github.com/coreos/fleet/Godeps/_workspace/src/golang.org/x/crypto/ssh"
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

// Creates a new SSH session
// and attaches a PTY to it.
func NewSession(config *Config, server string) (*gossh.Session, error) {
	client, err := createClient(config, server)

	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	
	if err != nil {
		return nil, err
	}
	
	modes := gossh.TerminalModes{
		gossh.ECHO:          0,     // disable echoing
		gossh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		gossh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	session.RequestPty("xterm", 80, 40, modes)

	return session, nil
}

// Creates a new SSH client based on the
// given configuration.
func createClient(config *Config, server string) (*ssh.SSHForwardingClient, error) {
	hostfile := ssh.NewHostKeyFile(config.Hostfile)
	checker := ssh.NewHostKeyChecker(hostfile)

	if config.Tunnel != "" {
		return ssh.NewTunnelledSSHClient(config.User, config.Tunnel, config.Address, checker, true, config.Timeout)
	}

	return ssh.NewSSHClient(config.User, config.Address, checker, true, config.Timeout)
}
