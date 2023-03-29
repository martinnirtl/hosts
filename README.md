# HOSTS CLI

Manage host entries - a simple tool for a simple, but annoying task!

```bash
Manage address mappings to ssh-config and hosts file. 
  Makes your life easier!
    Don't forget the sudo!

Usage:
  hosts [flags]
  hosts [command]

Available Commands:
  add         Add address mappings to ssh-config and hosts file
  completion  Generate completion script
  help        Help about any command
  print       Print contents of ssh-config and hosts file
  rm          Remove host entries from ssh-config and hosts file

Flags:
      --dry-run /etc/hosts   Only print updated /etc/hosts and `~/.ssh/config` files
  -h, --help                 help for hosts
      --hosts-file string    Set host file (e.g. ~/hosts). Default: /etc/hosts
      --ssh-config string    Set SSH Config file (e.g. /etc/ssh/config). Default: ~/.ssh/config

Use "hosts [command] --help" for more information about a command.
```

## Install

Installation is currently only supported via go:

```bash
go install github.com/martinnirtl/hosts
```
