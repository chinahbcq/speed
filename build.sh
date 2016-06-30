export GOPATH=`pwd`
go build -o bin/server src/main.go
go build -o bin/client src/test/test.go
