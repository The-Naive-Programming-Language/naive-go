GOPATH := $(shell go env GOPATH)
LINTER := $(GOPATH)/bin/staticcheck
GOFMT  := gofmt

SRC := $(shell find . -type f -name '*.go')

ENVS := GOOS=linux GOARCH=amd64

.PHONY: all run ci checkfmt fmt lint test

all: fmt lint test
	[[ ! -d build ]] && mkdir build || true
	go build -o build/naive .

ci: checkfmt lint test

checkfmt:
	$(GOFMT) -d $(SRC) || echo "use 'make fmt' to fix formatting issues"

fmt:
	$(GOFMT) -s -l -w $(SRC)

lint:
	# ignore the return code
	$(LINTER) ./... || true

test:
	$(ENVS) go test -coverprofile=/tmp/test.cov ./...

run: all
	build/naive
