#!/bin/bash

## CUSTOM BUILD VARS ##
OS=linux
ARCH=amd64
#######################

[ -z $1 ] || { 
    OS=$1
}

[ -z $2 ] || {
    OS=$2
}

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

BIN=./bin/tokenTools

[ "$OS" == "windows" ] && {
    BIN=${BIN}.exe
}

(
    env GOOS=$OS GOARCH=$ARCH \
        go build \
            -v \
            -o ${BIN} \
            -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT" \
            ./cmd/tokentools
)