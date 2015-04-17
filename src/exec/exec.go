// This package is used to "phisically"
// execute remote commands and handle their
// output on the current session.
package exec

import (
	golog "log"
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
		c := &ssh.Config{Address: serverConfig.Address, Alias: server, Tunnel: serverConfig.Tunnel, User: serverConfig.User, Hostfile: cfg.Hostfile}
		c.Timeout = time.Duration(10) * time.Second

		go func() {
			defer wg.Done()
			ExecuteRemoteCommand(command, c)
		}()
	}

	wg.Wait()
}

// Executes the given command through SSH,
// connecting with the given config.
func ExecuteRemoteCommand(command string, c *ssh.Config) {
	session, err := ssh.CreateSession(c)

	if err != nil {
		golog.Printf("Failed to create SSH session on %s", c.Address)
		return
	}

	stdout, stderr := log.GetRemoteLoggers(c.Alias)
	session.Stdout = stdout
	session.Stderr = stderr

	session.Run(command)
}
