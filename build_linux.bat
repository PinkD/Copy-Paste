@echo off

set GOOS=linux
set GOARCH=amd64
go build test.go
go build main.go
