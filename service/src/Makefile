export GOPROXY=https://mirrors.aliyun.com/goproxy/
export GO111MODULE=on

DATETIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GOFILES=$(shell find . -name "*.go" -type f -not -path "./vendor/*")
MODFILE=go.mod
VERSION=1.0.1

init:
	@if [ -f ${MODFILE} ] ; then rm ${MODFILE} ; fi
	@go mod init
	echo "package main\n\nconst version = \"${VERSION}\"\nconst tag = \"${DATETIME}\"\n" > ./cmd/version.go

build:
	go build -o blog  ./cmd
	upx blog

linux:
	GOOS=linux GOARCH=amd64 go build  -ldflags '-s -w' -o blog  ./cmd
	upx blog

