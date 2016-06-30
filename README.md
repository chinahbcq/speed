# Speed
Speed is a high performance framework designed for backend server.The default data protocol is NsHead and Json. The NsHead defined as follows:
```golang
type NsHead struct {
	id       uint16
	version  uint16
	LogId    uint32
	Provider string // max length 16
	MagicNum uint32
	reserved uint32
	bodyLen  uint32
}
```

# Installation
To install this package, you need to install Go and setup your Go workspace on your computer. The simplest way to install the package is to run:
```bash
$ go get -u github.com/chinahbcq/speed
```

# Prerequisites
This requires Go 1.6 or later.

# Quick Start

```bash
sh build.sh
nohup ./bin/server &
nohup ./bin/client &
```
