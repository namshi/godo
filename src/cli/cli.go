package cli

import (
	"fmt"
	"os"
	"strings"

	"./../config"
	"./../exec"
	"github.com/codegangsta/cli"
	"github.com/mgutz/ansi"
)

// Creates a new cli app.
func newApp() *cli.App {
	app := cli.NewApp()

	app.Name = "godo"
	app.Usage = "Stop SSHing into your server and run the same old commands. Automate. Automate. Automate."
	app.Version = "unstable"

	return app
}

// Registers all available commands
// on the app.
func addCommands(app *cli.App) {
	app.Action = func(c *cli.Context) {
		configFile := "./godo.yml"

		if c.String("config") != configFile {
			configFile = c.String("config")
		}

		cfg := config.Parse(configFile)
		cmd := c.Args().First()

		if command, ok := cfg.Commands[cmd]; ok {
			fmt.Printf("Executing '%s'", colorize(cmd))
			runCommand(command, cfg)
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
	for name, command := range commands {
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

// Colorizes a string for outputting on the CLI
func colorize(message string) string {
	return ansi.Color(message, "blue+h")
}

// Runs one of the commands stored in the config
// file.
//
// Here we will try to figure out if we need to execute
// it on a single server or a group of servers, and
// schedule the commands.
func runCommand(command config.Command, cfg config.Config) {
	fmt.Printf("\nCommand: '%s'", colorize(command.Exec))
	m := make(map[string]config.Server)

	if group, ok := cfg.Groups[command.Target]; ok {
		fmt.Printf("\nExecuting on servers %s", colorize(strings.Join(group, ", ")))

		for _, server := range group {
			m[server] = cfg.Servers[server]
		}
	} else {
		fmt.Printf("\nExecuting on server %s", colorize(command.Target))
		m[command.Target] = cfg.Servers[command.Target]
	}

	fmt.Println()
	fmt.Println()
	exec.ExecuteRemoteCommands(command.Exec, m, cfg)
}

// Adds global flags to the CLI app.
func addFlags(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Value:  "./godo.yml",
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
