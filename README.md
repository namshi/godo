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
and simply run the executable `godo`.

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

### Compiling

Alternatively, you can run and compile godo on
your own machine with a simple `go build -o godo main.go`.

At the same time, we provide a simple docker container
to run and compile it so that you don't have to
go to crazy if you don't have Go running on your
system:

```
git clone https://github.com/namshi/godo.git

cd godo

docker-compose run web go build -o godo main.go

./godo
```

The above command will compile the `godo` executable
in the current folder.

## Usage

Create a `godo.yml` file and put it in your home directory:

``` godo.yml
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
```

There are a few sections to keep in mind:

* `servers`: this is a dictionary of servers on which you can run commands on
* `groups`: a list of grouped servers (ie. you might want to group by role, AWS zone, etc)
* `commands`: the commands are the actual remote commands you would be executing on the servers,
they have a target (which can be a server or a group), the command that you would execute (`exec`)
and an optional description (which is printed when you do `godo help` or `godo`)
* `hostfile`: you can omit it, it's used not to always ask you to trust SSH host

## Tests

Run the tests with:

```
fig run godo go test ./...
```

## Todo

* check what happens if we dont provide any hostfile
* remove TODO
* custom SSH timeout
* implement @ operator to run command on a specific server
* add releases for a few platforms
* do not require ssh tunneling
* @all
* tests
  * config parsing
* autocomplete
* dont panic if you cant connect to via SSH or cant resolve an address