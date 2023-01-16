GOPATH := $(shell go env GOPATH)
LINTER := $(GOPATH)/bin/staticcheck
GOFMT  := gofmt

SRC := $(shell find . -type f -name '*.go')

ENVS := GOOS=linux GOARCH=amd64

.PHONY: all lint fmt test run

all: lint fmt test
	[[ ! -d build ]] && mkdir build || true
	go build -o build/naive .

lint:
	# ignore the return code
	$(LINTER) ./... || true

fmt:
	$(GOFMT) -s -l -w $(SRC)

test:
	# TODO
	$(ENVS) go test -coverprofile=/tmp/test.cov ./...

run: all
	build/naive
