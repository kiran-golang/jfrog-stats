#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

CGO_ENABLED=0
GOARCH=amd64
GOOS=linux

go build -o jfrog-stats -v main.go
go test -v -cover ./...