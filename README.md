# go-simple-memdb [![GoDoc](https://godoc.org/github.com/utrack/go-simple-memdb?status.svg)](https://godoc.org/github.com/utrack/go-simple-memdb)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/utrack/go-simple-memdb/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/utrack/go-simple-memdb.svg)](https://travis-ci.org/utrack/go-simple-memdb)
[![codecov.io](https://codecov.io/github/utrack/go-simple-memdb/coverage.svg?branch=master)](https://codecov.io/github/utrack/go-simple-memdb?branch=master)
[![Go Report Card](http://goreportcard.com/badge/utrack/go-simple-memdb)](http://goreportcard.com/report/utrack/go-simple-memdb)

[![forthebadge](http://forthebadge.com/images/badges/made-with-crayons.svg)](http://forthebadge.com)

Simple in-memory database written in Go.

# Requirements
Golang compiler and tools (v1.5 or later) are required. See the [official Getting Started guide](https://golang.org/doc/install) or your distro's docs for detailed instructions.

# Installation
```
go get -u github.com/utrack/go-simple-memdb
```

# Running
Check that your `PATH` envvar has `$GOPATH\bin` and run the command:
```
go-simple-memdb
```

# Protocol definition

## Data
* `SET <name <value>` – Sets the variable `name` to the value `value`. Variable name should not contain spaces.
* `GET <name>` – Value of the variable `name` is returned. `NULL` is returned if that variable was not set before.
* `UNSET <name>` – Unsets the variable name, making it just like that variable was never set.
* `NUMEQUALTO <value>` – Number of variables that are currently set to value is returned.
* `END` – Exit the program.

## Transactions
This storage supports nested transactions.
* `BEGIN` – Open a new transaction block. Transaction blocks can be nested; a `BEGIN` can be issued inside of an existing block.
* `ROLLBACK` – Most recent transaction block is closed, all changes in it are forgotten. `NO TRANSACTION` is printed if there's no transactions in progress.
* `COMMIT` – Closes all open transaction blocks, permanently applying the changes made in them. `NO TRANSACTION` is printed if there's no transactions in progress.

Any data command that is run outside of a transaction block is committed immediately.

# Testing
```
go test github.com/utrack/go-simple-memdb/...
```
Tests are using the [GoConvey](https://github.com/smartystreets/goconvey) framework. If you have `goconvey` tools installed in your `$PATH`, run `goconvey github.com/utrack/go-simple-memdb/...` to use its web interface.
