// This package is used to "phisically"
// execute remote commands and handle their
// output on the current session.
package exec

import (
	"fmt"
	"github.com/mgutz/ansi"
	goexec "os/exec"
	"strings"
	"sync"
	"time"

	"github.com/namshi/godo/src/config"
	"github.com/namshi/godo/src/log"
	"github.com/namshi/godo/src/ssh"
)

// Executes the command on the given
// servers.
//
// The main goal of this method is to
// figure out whether this command needs
// to be executed locally or remotely,
// and go ahead with the proper execution
// strategy.
func ExecuteCommands(command string, servers map[string]config.Server, cfg config.Config) {
	var wg sync.WaitGroup
	wg.Add(len(servers))

	for server, serverConfig := range servers {
		if server == "local" {
			go executeLocalCommand(command, wg.Done)
		} else {
			go executeRemoteCommand(command, server, serverConfig, cfg, wg.Done)
		}
	}

	wg.Wait()
}

// Parses a command in format that
// is suitable for exec.Command().
//
// In practice, ls -la /tmp becomes
// "ls" and ["-la", "/tmp"].
func parseLocalCommand(command string) (string, []string) {
	args := strings.Fields(command)

	return args[0], args[1:len(args)]
}

// Executes a command locally.
func executeLocalCommand(command string, done func()) {
	c, args := parseLocalCommand(command)
	cmd := goexec.Command(c, args...)
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

	defer done()
}

// Executes the given command through SSH,
// connecting with the given config.
func executeRemoteCommand(command string, server string, serverConfig config.Server, cfg config.Config, done func()) {
	c := &ssh.Config{Address: serverConfig.Address, Alias: server, Tunnel: serverConfig.Tunnel, User: serverConfig.User, Hostfile: cfg.Hostfile}
	c.Timeout = time.Duration(cfg.Timeout) * time.Second
	session, _ := ssh.CreateClient(c).NewSession()

	stdout, stderr := log.GetRemoteLoggers(server)
	session.Stdout = stdout
	session.Stderr = stderr

	session.Run(command)
	defer done()
}
