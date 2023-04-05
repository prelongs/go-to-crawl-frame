#!/bin/bash
APP_NAME="task_vod"
export GOOS=linux GOARCH=amd64
go build -o ./build/${APP_NAME} main.go
#upx  ./build/${APP_NAME}

export GOOS=windows
go build -o ./build/${APP_NAME}.exe main.go



#APP_NAME="crawl_down"
#export GOOS=linux GOARCH=amd64
#go build
#go build -o ./build/${APP_NAME} -ldflags="-w -s" crawl_down.go
#upx  ./build/${APP_NAME}

