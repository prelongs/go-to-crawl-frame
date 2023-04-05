#!/bin/bash
APP_NAME="task"
export GOOS=linux GOARCH=amd64
go build -o ./build/${APP_NAME} -ldflags="-w -s" main.go
upx  ./build/${APP_NAME}

#APP_NAME="crawl_down"
#export GOOS=linux GOARCH=amd64
#go build -o ./build/${APP_NAME} -ldflags="-w -s" crawl_down.go
#upx  ./build/${APP_NAME}

