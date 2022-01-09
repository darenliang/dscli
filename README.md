# dscli: A command-line tool for storing files on Discord

[![Go Report Card](https://goreportcard.com/badge/github.com/darenliang/dscli)](https://goreportcard.com/report/github.com/darenliang/dscli)
[![License](https://img.shields.io/github/license/darenliang/dscli)](https://github.com/darenliang/dscli/blob/master/LICENSE)

Dscli (Discord store CLI) provides a way to store files with no size restrictions.

## Requirements

* Empty Discord server
* You'll need a bot/user token and the server id

## Installation

Make sure you have go installed (version 1.13+ is required).

```
go get -u github.com/darenliang/dscli
```

You can also download pre-built binaries for Windows, Linux and MacOS: https://github.com/darenliang/dscli/releases

## Quickstart

For Bot Tokens (highly recommended):

```
dscli config -t=<YOUR-DISCORD-BOT-TOKEN> -i=<YOUR-SERVER-ID> -b
```

For User Tokens:

```
dscli config -t=<YOUR-DISCORD-TOKEN> -i=<YOUR-SERVER-ID>
```

A complete setup guide (including creating a server and getting user tokens) can be found [here](https://github.com/darenliang/dscli/blob/master/quickstart/README.md) or by using the command:

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

#### Statistics

* Number of files stored
* Total size stored
* Upload chunk size

```
dscli stats
```

## Quick Benchmarks

Ubuntu 18.04.5 LTS on Gigabit internet

Please take these numbers with a grain of salt.

|                     | 10 MB         | 100 MB        | 1000 MB       |
| ------------------- | ------------- | ------------- | ------------- |
| Upload              | 18.161 MB/s   | 7.760 MB/s    | 6.352 MB/s    |
| Upload Nitro        | 22.450 MB/s   | 22.640 MB/s   | 20.381 MB/s   |
| Download Cold       | 23.878 MB/s   | 32.895 MB/s   | 32.976 MB/s   |
| Download Cold Nitro | 22.586 MB/s   | 39.553 MB/s   | 34.811 MB/s   |
| Download Warm       | 57.061 MB/s   | 74.972 MB/s   | 75.638 MB/s   |
| Download Warm Nitro | 65.748 MB/s   | 78.645 MB/s   | 69.352 MB/s   |

* Cold: Initial download
* Warm: Later downloads
* Nitro: Discord Nitro

## Limitations

* No folders
* Limited filename lengths
    * Filenames are encoded in base32
    * Base32 strings cannot be longer than 100 characters
* No file appending or editing
* File number limit of 500 (max number of channels)
* Subject to Discord rate limits (5 anything per 5 seconds per server)
    * 8 MB x 5 / 5s = 8 MB/s (theoretical sustained upload speed for non-Nitro users)
    * 50 MB x 5 / 5s = 50 MB/s (theoretical sustained upload speed for Nitro classic users)
    * 100 MB x 5 / 5s = 100 MB/s (theoretical sustained upload speed for Nitro users)

## License

[MIT](https://github.com/darenliang/dscli/blob/master/LICENSE)
