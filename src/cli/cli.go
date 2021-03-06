package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/namshi/godo/src/config"
	"github.com/namshi/godo/src/exec"
	"github.com/namshi/godo/src/log"
)

const defaultConfigFile = "./godo.yml"

// Creates a new cli app.
func newApp() *cli.App {
	app := cli.NewApp()

	app.Name = "godo"
	app.Usage = "Stop SSHing into your server and run the same old commands. Automate. Automate. Automate."
	app.Version = "v1.3.0"

	return app
}

// Parses the arguments that we got
// from the command line.
//
// Arguments have the format command[@target],
// so something like "uptime" or "uptime @ load-balancer".
func parseArgs(c *cli.Context) (string, string) {
	cmds := strings.Split(strings.Replace(strings.Join(c.Args(), ""), " ", "", -1), "@")

	if len(cmds) == 2 {
		return cmds[0], cmds[1]
	}

	return cmds[0], ""
}

// Registers all available commands
// on the app.
func addCommands(app *cli.App) {
	app.Action = func(c *cli.Context) {
		configFile := defaultConfigFile

		if c.String("config") != configFile {
			configFile = c.String("config")
		}

		cfg := config.Parse(configFile)
		cmd, target := parseArgs(c)

		if command, ok := cfg.Commands[cmd]; ok {
			log.Info("Executing '%s'", cmd)

			if target == "" {
				target = command.Target
			}

			runCommand(command, cfg, target)
		} else {
			printAvailableCommands(app, cfg.Commands, c)
		}
	}
}

// Helper function that prints all available
// commands.
//
// @todo we should register commands before this,
// but somehow if I do it the cli app executes a
// random command
func printAvailableCommands(app *cli.App, commands map[string]config.Command, c *cli.Context) {
	// golang's map iteration is random but we want commands to be printed in alphabetical order
	// @see http://nathanleclaire.com/blog/2014/04/27/a-surprising-feature-of-golang-that-colored-me-impressed/
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		command := commands[name]
		description := command.Exec

		if command.Description != "" {
			description = command.Description
		}

		cmd := cli.Command{
			Name:  name,
			Usage: description,
		}

		app.Commands = append(app.Commands, cmd)
	}

	app.Command("help").Run(c)
}

// Adds servers to the list of targets
// by checking if the specified
// target was a group: if it was, add
// all servers in that group.
func addTargetFromGroups(targets map[string]config.Server, target string, cfg config.Config) {
	if group, ok := cfg.Groups[target]; ok {
		for _, server := range group {
			targets[server] = cfg.Servers[server]
		}
	}
}

// Adds a server to the list of targets if
// the specified target was a single server.
func addTargetFromServer(targets map[string]config.Server, target string, cfg config.Config) {
	if _, ok := cfg.Servers[target]; ok {
		targets[target] = cfg.Servers[target]
	}
}

// Returns a list of targets on which we
// should execute the given command, based
// on the target specified by the user.
//
// If the user specifies a target, we use that
// one; if there is no user-specified target
// we simply look at the configuration of the
// command.
//
// A target can be a server, group of servers
// or a special alias.
//
// The supported aliases are
// - all: will execute the command on all servers
// - local: instead of executing the command remotely
//          it will execute it on the current machine
func getTargets(command config.Command, cfg config.Config, target string) map[string]config.Server {
	targets := make(map[string]config.Server)

	for _, target = range strings.Split(target, ",") {
		if target == "all" {
			targets = cfg.Servers
		} else if target == "local" {
			targets["local"] = config.Server{}
		} else {
			addTargetFromGroups(targets, target, cfg)
			addTargetFromServer(targets, target, cfg)
		}
	}

	return targets
}

// Runs one of the commands stored in the config
// file.
func runCommand(command config.Command, cfg config.Config, target string) {
	log.Info("\nCommand: '%s'", command.Exec)
	targets := getTargets(command, cfg, target)
	targetNames := []string{}

	for serverName, _ := range targets {
		targetNames = append(targetNames, serverName)
	}

	if len(targets) > 0 {
		log.Info("\nExecuting on server '%s'", strings.Join(targetNames, ", "))
		fmt.Println()
		fmt.Println()
		exec.ExecuteCommands(command.Exec, targets, cfg)
	} else {
		log.Err("\nNo target server / group with the name '%s' could be found, maybe a typo?", target)
		printAvailableTargets(cfg)
	}
}

// Outputs the available targets,
// extracted from the config file.
func printAvailableTargets(cfg config.Config) {
	log.Err("\n\nAvailable groups are:")

	for groupName, group := range cfg.Groups {
		log.Err("\n  * %s (%s)", groupName, strings.Join(group, ", "))
	}

	log.Err("\n\nAvailable servers are:")

	for serverName := range cfg.Servers {
		log.Err("\n  * %s", serverName)
	}
}

// Adds global flags to the CLI app.
func addFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  defaultConfigFile,
			Usage:  "configuration file to be used for running godo",
			EnvVar: "GODO_CONFIG",
		},
	}
}

// Runs the cli app!
func Run() {
	app := newApp()
	addCommands(app)
	addFlags(app)

	app.Run(os.Args)
}
