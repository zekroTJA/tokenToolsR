#!/bin/bash

DATE=$(date -u '+%Y-%m-%d_%I:%M:%S%p')
TAG=$(git describe --tags)
COMMIT=$(git rev-parse HEAD)

(env GOOS=linux GOARCH=amd64 go build -v -o tokenTools -ldflags "-X main.appDate=$DATE -X main.appVersion=$TAG -X main.appCommit=$COMMIT")