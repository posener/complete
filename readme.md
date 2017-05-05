# complete

[![Build Status](https://travis-ci.org/posener/complete.svg?branch=master)](https://travis-ci.org/posener/complete)
[![codecov](https://codecov.io/gh/posener/complete/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/complete)

WIP

a tool for bash writing bash completion in go.

## example: `go` command bash completion

Install in you home directory:

```
go build -o ~/.bash_completion/go ./gocomplete
echo "complete -C ~/.bash_completion/go go" >> ~/.bashrc
```

Or, install in the root directory:

```
sudo go build -o /etc/bash_completion.d/go ./gocomplete
```

