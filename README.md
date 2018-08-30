# zeno
[![Go Report Card](https://goreportcard.com/badge/github.com/meomap/zeno)](https://goreportcard.com/report/github.com/meomap/zeno) [![Build Status](https://travis-ci.org/meomap/zeno.svg?branch=master)](https://travis-ci.org/meomap/zeno) [![Coverage Status](https://coveralls.io/repos/github/meomap/zeno/badge.svg?branch=master)](https://coveralls.io/github/meomap/zeno?branch=master)

CLI tool to examine ansible playbooks affected by git changes

## Installation
```
$ go get -u github.com/meomap/zeno
```
## Usage
```
$ zeno -files="$(git diff $COMMIT_HASH_BEFORE $COMMIT_HASH_AFTER --name-only)" -playbooks=qa/site.yml,staging/site.yml
qa/site.yml,staging/site.yml
```
## Features

- Ansible playbook supported.

## Contributing

Bug reports & pull requests are welcome.
