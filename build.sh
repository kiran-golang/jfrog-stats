#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=0
export GOARCH=amd64
export GOOS=linux

go build -o jfrog-stats -v main.go
go test -v -cover ./...