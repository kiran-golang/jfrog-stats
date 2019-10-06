BINARY := jfrog-stats
PLATFORM := linux

export GO111MODULE=on

all: test build

build: clean
	GOARCH=amd64 GOOS=$(PLATFORM) go build -mod vendor -o $(BINARY) -v main.go

clean:
	go clean
	rm -f $(BINARY)

test:
	go test -cover ./...

format:
	go fmt ./...

.PHONY: test