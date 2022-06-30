GO ?= go

all: tests little-vm-helper 

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	$(GO) build

FORCE:
