export ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export PKG := github.com/zimmski/osutil

export UNIT_TEST_TIMEOUT := 480

ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(ARGS):;@:) # turn arguments into do-nothing targets
export ARGS

all: dependencies install test-verbose
.PHONY: all

clean:
	go clean -i $(PKG)/...
	go clean -i -race $(PKG)/...
.PHONY: clean

dependencies:
	go get -t -v ./...
	go build -v ./...
.PHONY: dependencies

install:
	go install -v ./...
.PHONY: install

test:
	go test -race -test.timeout $(UNIT_TEST_TIMEOUT)s $(PKG_TEST)
.PHONY: test

test-verbose:
	go test -race -test.timeout $(UNIT_TEST_TIMEOUT)s -v $(PKG_TEST)
.PHONY: test-verbose
