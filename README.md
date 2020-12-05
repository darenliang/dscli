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

|                 | 10 MB         | 100 MB        | 1000 MB       |
| --------------- | ------------- | ------------- | ------------- |
| Upload          | 18.161 MB/s   | 7.760 MB/s    | 6.352 MB/s    |
| Download (cold) | 23.878 MB/s   | 32.895 MB/s   | 32.976 MB/s   |
| Download (warm) | 57.061 MB/s   | 74.972 MB/s   | 75.638 MB/s   |

* cold: Initial download
* warm: Later downloads

## Limitations

* No folders
* Limited filename lengths
    * Filenames are encoded in base32
    * Base32 strings cannot be longer than 100 characters
* No file appending or editing
* File number limit of 500 (max number of channels)
* Subject to Discord rate limits (5 anything per 5 seconds per server)
    * 8 MB x 5 / 5s = 8 MB/s (theoretical sustained upload speed)

## License

[MIT](https://github.com/darenliang/dscli/blob/master/LICENSE)
