# dscli: A command-line tool for storing files on Discord

[![Go Report Card](https://goreportcard.com/badge/github.com/darenliang/dscli)](https://goreportcard.com/report/github.com/darenliang/dscli)
[![License](https://img.shields.io/github/license/darenliang/dscli)](https://github.com/nikel-api/nikel/blob/master/LICENSE)

Dscli stands for Discord store CLI and provides a way to store files with no size restrictions.

### Installation

Make sure you have go installed (version 1.13+ is required).

```
go get -u github.com/darenliang/dscli
```

### Quickstart

```
dscli quickstart
```

### Examples

Configure application

```
dscli config
```

List files
```
dscli ls

test.txt
test.mp3
test.mp4
```

Upload file
```
dscli up foo.txt

Uploading foo.txt 100% |███████| (1/1 MB, 2.5 MB/s)
```

Download file
```
$ dscli dl foo.txt

Downloading bar.txt 100% |███████| (1/1 MB, 5.0 MB/s)
```

Move file
```
dscli mv foo.txt bar.txt
```

Remove file
```
dscli rm foo.txt
```

### Limitations

* No folders
* Limited filename lengths
* No file appending or editing
* File number limit of 500
* Subject to Discord rate limits

### License

[MIT](https://github.com/darenliang/dscli/blob/master/LICENSE)