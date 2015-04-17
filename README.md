# godo

Remote execution level 9000: **go** and **do**
stuff.

`godo` is a very simple yet powerful tool that
let's you specify a list of repetitive / useful
commands you run on remote hosts, and run them
with ease, without having to remember them or,
worse, login on each server and execute them
manually.

## Installation

Grab the [latest release](https://github.com/namshi/godo/releases)
and simply run the executable:

```
~/projects/namshi/godo (master ✘)✭ ᐅ ./godo
NAME:
   godo - Stop SSHing into your server and run the same old commands. Automate. Automate. Automate.

USAGE:
   godo [global options] command [command options] [arguments...]

VERSION:
   X.Y.Z

AUTHOR(S): 
   
COMMANDS:
   help, h	Shows a list of commands or help for one command
   uptime	Retrieves uptime info for all machines
   slow-queries	mysql-log -n 10
   nginx-logs	sudo tail -10f /var/log/nginx/access.log
   
GLOBAL OPTIONS:
   --help, -h			show help
   --version, -v		print the version
   
~/projects/namshi/godo (master ✘)✭ ᐅ ./godo uptime
Executing 'uptime'
Command: 'uptime'
Executing on servers web, db

...
...
```

## Usage

Create a `godo.yml` file and put it in your home directory:

``` yaml
servers:
  web:
    address: "xxx.xxx.xx.xxx:22"
    user:    "me"
  db:
    address: "xxx.xxx.xx.xxx:22"
    tunnel:  "tunnel.yourcompany.com:22"
    user:    "me"
groups:
  all: [web, db]
commands:
  uptime:
    target:       all
    exec:         "uptime"  
    description:  "Retrieves uptime info for all machines"  
  slow-queries:
    target: db
    exec: "mysql-log -n 10"
  nginx-logs:
    target: web
    exec: "sudo tail -10f /var/log/nginx/access.log"
hostfile: "/home/YOU/.ssh/known_hosts"
timeout: 2
```

There are a few sections to keep in mind:

* `servers`: this is a dictionary of servers on which you can run commands on
* `groups`: a list of grouped servers (ie. you might want to group by role, AWS zone, etc)
* `commands`: the commands are the actual remote commands you would be executing on the servers,
they have a target (which can be a server or a group), the command that you would execute (`exec`)
and an optional description (which is printed when you do `godo help` or `godo`)
* `hostfile`: you can omit it, it's used not to always ask you to trust SSH host

Godo will try to read the `godo.yml` configuration file
from 3 different directories:

* your home
* the current directory
* the directory from which the godo executable runs

but you can also specify the path to a different
configuration file with the `-c` or `--config` flags:

```
godo -c ./../my-config.yml mysql-log
```

Sometimes, though, you might want to run a command
that you usually execute on some servers on a
different server, and you can do it by simply
specifying it from the command line:

```
godo uptime @ db

# or

godo uptime@db
```

Godo provides a special group, called `all`, that
represents all servers, so you can always run
something like `godo uptime @ all`.

## Additional documentation

You can run the docs through `godoc -http=:6060 -path=.`.

## Gotchas

Currently all servers need to be in your `known_hosts` file (ie. you
have to have SSHed into them at least once before using them with godo).

## Compiling

Alternatively, you can run and compile godo on
your own machine with a simple `go build -o godo main.go`.

At the same time, we provide a simple docker container
to run and compile it so that you don't have to
go to crazy if you don't have Go running on your
system:

```
git clone https://github.com/namshi/godo.git

cd godo

docker-compose run godo gox --output=build/{{.OS}}_{{.Arch}}/{{.Dir}}

./godo
```

The above command will compile the `godo` executables
(for various platforms) in the `build` folder.

We use a simple makefile to create new releases and
yoou can probably do the same: just run `make` in the
root of the repo and check the `build` folder. This
requires that `docker` and `docker-compose` are installed
on your system.
