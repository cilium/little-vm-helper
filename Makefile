GO ?= go

all: tests little-vm-helper 

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	$(GO) build ./cmd/lvh

clean:
	rm -f lvh
FORCE:
