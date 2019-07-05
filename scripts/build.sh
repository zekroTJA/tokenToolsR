#!/bin/bash

## CUSTOM BUILD VARS ##
OS=linux
ARCH=amd64
#######################

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

(env GOOS=$OS GOARCH=$ARCH \
    go build \
        -v \
        -o tokenTools \
        -ldflags "\
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppDate=$DATE \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppVersion=$TAG \
            -X github.com/zekroTJA/tokenToolsR/internal/static.AppCommit=$COMMIT" \
        ./cmd/tokentools/*.go)