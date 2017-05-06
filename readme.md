# complete

[![Build Status](https://travis-ci.org/posener/complete.svg?branch=master)](https://travis-ci.org/posener/complete)
[![codecov](https://codecov.io/gh/posener/complete/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/complete)

WIP

A tool for bash writing bash completion in go.

Writing bash completion scripts is a hard work. This package provides an easy way
to create bash completion scripts for any command, and also an easy way to install/uninstall
the completion of the command.

## go command bash completion

In [gocomplete](./gocomplete) there is an example for bash completion for the `go` command line.

### Install

```
go get github.com/posener/complete/gocomplete
gocomplete -install
```

### Uninstall

```
gocomplete -uninstall
```
