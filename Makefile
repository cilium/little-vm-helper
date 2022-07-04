GO ?= go

all: tests little-vm-helper 

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	$(GO) build ./cmd/little-vm-helper

clean:
	rm -f little-vm-helper
FORCE:
