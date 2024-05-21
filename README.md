# GitBatch Script

Script to batch handling git projects

# Installation

Require go 1.22+

```shell
go install github.com/sitdownrightnow2552/gitbatch@latest
```

## Feature

By default, this script is configured to use ssh auth. To using basic auth, you must specify your username
using ``--user=<username>`` or ``-u=<username>``. To switch back to ssh mode, specify ```--user=@ssh```

To show list of available commands

```shell
gitbatch help
```

### Clone all project in gitlab group

```shell
gitbatch clone [group id] [dir?]
```

```shell
gitbatch cg [group id] [dir?]
```

### Fetch all project inside directory

```shell
gitbatch fetch [dir?]
```

Known issue: if you're using @ssh mode then your may see SSH_AUTH_SOCK error. To fix this issue consider lowering
the ``--parallel`` value. Currently, the parallel is capped at 8 for this command only.

## MacOS User

Using @ssh mode you may need to add key manually:

```shell
chmod 600 ~/.ssh/id_rsa
ssh-add ~/.ssh/id_rsa
```

Reference: [Go Git Issue: ssh private key not being picked up](https://github.com/go-git/go-git/issues/218)