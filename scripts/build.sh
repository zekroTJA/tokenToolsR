#!/bin/bash

## CUSTOM BUILD VARS ##
OS=linux
ARCH=amd64
#######################

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

(env GOOS=$OS GOARCH=$ARCH go build -v -o tokenTools -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")