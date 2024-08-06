# Git Batch

Script to batch handling git projects

# Installation

Require go 1.22+

```shell
go install github.com/mawngo/gitbatch@latest
```

## Feature

By default, this script is configured to use ssh auth. To using basic auth, you must specify your username
using ``--user=<username>`` or ``-u=<username>``. To switch back to ssh mode, specify ```--user=@ssh```

To show list of available commands

```
> gitbatch -h

Apply git command to all sub folder

Usage:
  gitbatch [command]

Available Commands:
  clone       Clone all project in group (alias for 'clone gitlab')
  fetch       Fetch all project in directory
  pull        Pull all project in directory
  push        Push all project in directory
  help        Help about any command
  completion  Generate the autocompletion script for the specified shell

Flags:
  -m, --mode string    Host mode (default "gitlab")
      --parallel int   Maximum parallel for each commands (default 32)
      --token string   Host token
  -u, --user string    Auth user name [<user>, @ssh] (default "@ssh")
  -h, --help           help for gitbatch

Use "gitbatch [command] --help" for more information about a command.
```

### Clone all projects in gitlab group

```
gitbatch clone [group id] [dir?]
```

```
gitbatch cg [group id] [dir?]
```

### Fetch all projects inside directory

```
gitbatch fetch [dir?]
```

Known issue: if you're using @ssh mode then your may see SSH_AUTH_SOCK error. To fix this issue consider lowering
the ``--parallel`` value. Currently, the parallel is capped at 8 for this command only.

## MacOS User

Using @ssh mode, you may need to add key manually:

```shell
chmod 600 ~/.ssh/id_rsa
ssh-add ~/.ssh/id_rsa
```

Reference: [Go Git Issue: ssh private key not being picked up](https://github.com/go-git/go-git/issues/218)