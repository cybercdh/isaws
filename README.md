# isaws

Takes a list of IP Addresses from stdin and outputs those which appear in current AWS IP Ranges

## Install

If you have Go installed and configured (i.e. with `$GOPATH/bin` in your `$PATH`):

```
go install github.com/cybercdh/isaws@latest
```

## Usage

```
$ echo 1.2.3.4 | isaws
```
or 
```
$ cat <file> | isaws
```

If an IP is found to be within the current list of AWS IP ranges, then details will be printed

### Options

```
Usage of isaws:
  -c int
    	set the concurrency level (default 50)
```
