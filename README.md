# zeno
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
