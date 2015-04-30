// This package is used to "phisically"
// execute remote commands and handle their
// output on the current session.
package exec

import (
	"fmt"
	gossh "github.com/coreos/fleet/Godeps/_workspace/src/golang.org/x/crypto/ssh"
	"github.com/mgutz/ansi"
	goexec "os/exec"
	"sync"
	"time"

	"github.com/namshi/godo/src/config"
	"github.com/namshi/godo/src/log"
	"github.com/namshi/godo/src/ssh"
)

// Checks whether a command needs
// to be run locally or remotely.
//
// If the only server is "local"
// then it means that we need to
// simply run the command on this
// local machine.
func isLocalCommand(servers map[string]config.Server) bool {
	if len(servers) == 1 {
		if _, ok := servers["local"]; ok {
			return true
		}
	}

	return false
}

// Executes the command on the given
// servers.
//
// The main goal of this method is to
// figure out whether this command needs
// to be executed locally or remotely,
// and go ahead with the proper execution
// strategy.
func ExecuteCommands(command string, servers map[string]config.Server, cfg config.Config) {
	if isLocalCommand(servers) {
		executeLocalCommand(command)
	} else {
		executeRemoteCommands(command, servers, cfg)
	}
}

// Executes a command locally.
func executeLocalCommand(command string) {
	cmd := goexec.Command(command)
	stdout, stderr := log.GetRemoteLoggers("local")
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Start()

	// failed to spawn new process
	if err != nil {
		fmt.Println(ansi.Color(err.Error(), "red+h"))
	}

	// Failed to execute?
	err = cmd.Wait()
	if err != nil {
		fmt.Println(ansi.Color(err.Error(), "red+h"))
	}
}

// Executes the given command on a series
// of servers.
//
// We will launch N goroutins based on how
// many commands we need to remotely execute,
// and stop the execution once everyone is
// done.
func executeRemoteCommands(command string, servers map[string]config.Server, cfg config.Config) {
	var wg sync.WaitGroup
	wg.Add(len(servers))

	for server, serverConfig := range servers {
		go func(server string, serverConfig config.Server) {
			c := &ssh.Config{Address: serverConfig.Address, Alias: server, Tunnel: serverConfig.Tunnel, User: serverConfig.User, Hostfile: cfg.Hostfile}
			c.Timeout = time.Duration(cfg.Timeout) * time.Second
			session, _ := ssh.CreateClient(c).NewSession()
			executeRemoteCommand(command, session, server)
			defer wg.Done()
		}(server, serverConfig)
	}

	wg.Wait()
}

// Executes the given command through SSH,
// connecting with the given config.
func executeRemoteCommand(command string, session *gossh.Session, server string) {
	stdout, stderr := log.GetRemoteLoggers(server)
	session.Stdout = stdout
	session.Stderr = stderr

	session.Run(command)
}
