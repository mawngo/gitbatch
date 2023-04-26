# GitBatch Script

Script to batch handling git projects

# Installation

Require git installed

```shell
git clone https://github.com/lana-toolbox/gitbatch-script gitbatch-script
cd gitbatch-script
go install .
```

## Feature

By default, this script is configured to use ssh auth. To using basic auth, you must specify your username
using ``--user=<username>``. To switch back to ssh mode, specify ```--user=@ssh```

To show list of available commands

```shell
gitbatch help
```

### Clone all project in gitlab group

```shell
gitbatch clone [group id] [dir?]
```

### Fetch all project inside directory

```shell
gitbatch fetch [dir?]
```

Known issue: if you're using @ssh mode then your may see SSH_AUTH_SOCK error. To fix this issue consider lowering
the ``--parallel`` value. Currently, the parallel is capped at 8 for this command only.