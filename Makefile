.PHONY: all clean

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

BINARY=crawler
VERSION=0.1.0
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all: build docker

build:
	go build -o $(BINARY) $(LDFLAGS)

docker:
	docker build -t "javking07/servers/$(BINARY):$(VERSION)" \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

clean:
-rm $(BINARY)