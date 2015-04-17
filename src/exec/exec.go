// This package is used to "phisically"
// execute remote commands and handle their
// output on the current session.
package exec

import (
	gossh "github.com/coreos/fleet/Godeps/_workspace/src/golang.org/x/crypto/ssh"
	"sync"
	"time"

	"./../config"
	"./../log"
	"./../ssh"
)

// Executes the given command on a series
// of servers.
//
// We will launch N goroutins based on how
// many commands we need to remotely execute,
// and stop the execution once everyone is
// done.
func ExecuteRemoteCommands(command string, servers map[string]config.Server, cfg config.Config) {
	var wg sync.WaitGroup
	wg.Add(len(servers))

	for server, serverConfig := range servers {
		go func(server string, serverConfig config.Server) {
			c := &ssh.Config{Address: serverConfig.Address, Alias: server, Tunnel: serverConfig.Tunnel, User: serverConfig.User, Hostfile: cfg.Hostfile}
			c.Timeout = time.Duration(60) * time.Second
			session, _ := ssh.CreateClient(c).NewSession()
			ExecuteRemoteCommand(command, session, server)
			defer wg.Done()
		}(server, serverConfig)
	}

	wg.Wait()
}

// Executes the given command through SSH,
// connecting with the given config.
func ExecuteRemoteCommand(command string, session *gossh.Session, server string) {
	stdout, stderr := log.GetRemoteLoggers(server)
	session.Stdout = stdout
	session.Stderr = stderr

	session.Run(command)
}
