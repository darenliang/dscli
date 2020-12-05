# dscli: A command-line tool for storing files on Discord

[![Go Report Card](https://goreportcard.com/badge/github.com/darenliang/dscli)](https://goreportcard.com/report/github.com/darenliang/dscli)
[![License](https://img.shields.io/github/license/darenliang/dscli)](https://github.com/nikel-api/nikel/blob/master/LICENSE)

Dscli (Discord store CLI) provides a way to store files with no size restrictions.

## Discord Requirements

* Empty Discord server
* Invited Discord bot that has permissions to manage servers
* You'll need the Discord bot token and the server id

## Installation

Make sure you have go installed (version 1.13+ is required).
```
go get -u github.com/darenliang/dscli
```

## Quickstart

A complete setup guide can be found at:
```
dscli quickstart
```

## Commands

#### Configure application
```
dscli config
```

#### List files
```
dscli ls
```

#### Upload file
Use local filename
```
dscli up <local file>
```
Specify remote filename
```
dscli up <local file> <remote file>
```

#### Download file
Use remote filename
```
dscli dl <remote file>
```
Specify local filename
```
dscli dl <remote file> <local file>
```

#### Move file
```
dscli mv <source file> <destination file>
```

#### Remove file
```
dscli rm <remote file>
```

## Quick Benchmarks

Ubuntu 18.04.5 LTS on Gigabit internet

```
dscli up 100MB.bin

Uploading 100MB.bin 100% |███████████████████████████| (100/100 MB, 7.509 MB/s)

dscli dl 100MB.bin

Downloading 100MB.bin 100% |████████████████████████| (100/100 MB, 14.710 MB/s)
```

## Limitations

* No folders
* Limited filename lengths
* No file appending or editing
* File number limit of 500
* Subject to Discord rate limits

## License

[MIT](https://github.com/darenliang/dscli/blob/master/LICENSE)
