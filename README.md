# Hosts CLI

Manage host entries - a simple tool for a simple, but annoying task!

```bash
Manage address mappings of SSH config and optionally for hosts file.

Usage:
  hosts [flags]
  hosts [command]

Available Commands:
  add         Add address mappings to ssh-config and hosts file
  completion  Generate completion script
  edit        Edit host entries of SSH config and optionally hosts file
  help        Help about any command
  print       Print contents of ssh-config and hosts file
  rm          Remove one or more host entries from ssh-config and hosts file
  version     Print CLI version information

Flags:
      --dry-run             Only print updated /etc/hosts and ~/.ssh/config files
      --etc-hosts           Additionally add entry to /etc/hosts file (requires sudo)
  -h, --help                help for hosts
      --hosts-file string   Set host file (e.g. ~/hosts); default: /etc/hosts
      --ssh-config string   Set SSH Config file (e.g. /etc/ssh/config); default: ~/.ssh/config

Use "hosts [command] --help" for more information about a command.

```

## Install

The CLI is available via a Brew Tap. Run the following command to install the Hosts CLI

```bash
brew install martinnirtl/tap/hosts
```

### Go install

Alternatively, Hosts CLI can be installed from source via Go:

```bash
go install github.com/martinnirtl/hosts-cli
```

> Please note: Installation via Go installs the CLI bin as `hosts-cli`, not `hosts`.
