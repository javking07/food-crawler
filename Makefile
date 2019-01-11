.PHONY: all clean

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

GO_PKGS=$(shell go list ./... | grep -v -e "/scripts")
BINARY=food-crawler
VERSION=0.1.0
BUILD=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

all: check-gofmt test build docker

check-gofmt:
	@echo "Checking formatting..."
	@FMT="0"; \
	for pkg in $(GO_PKGS); do \
		OUTPUT=`gofmt -l $(GOPATH)/src/$$pkg/*.go`; \
		if [ -n "$$OUTPUT" ]; then \
			echo "$$OUTPUT"; \
			FMT="1"; \
		fi; \
	done ; \
	if [ "$$FMT" -eq "1" ]; then \
		echo "Problem with formatting in files above."; \
		exit 1; \
	else \
		echo "Success - way to run gofmt!"; \
	fi

test:
	go test -v $(GO_PKGS)

build:
	go build -o $(BINARY) $(LDFLAGS)

docker:
	docker build -t "javking07/$(BINARY):$(VERSION)" \
		--build-arg build=$(BUILD) --build-arg version=$(VERSION) \
		-f Dockerfile .

docker-test:
	@docker rmi -f $(BINARY) && docker build -t $(BINARY) . && docker run -v $(pwd):/app/config.json \
	 $(BINARY) --config=config.json

clean:
	rm $(BINARY)

echo:
	@echo $(pwd)