TARGETS := \
	github.com/modcloth-labs/hookworm \
	github.com/modcloth-labs/hookworm/hookworm-server

ADDR := :9988


all: test

test: build
	go test -v $(TARGETS)

build: deps
	go install -x $(TARGETS)

deps:
	go get -x $(TARGETS)

serve:
	$${GOPATH%%:*}/bin/hookworm-server -a $(ADDR) -S


.PHONY: all test build deps
