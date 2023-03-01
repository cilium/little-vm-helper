GO ?= go

OCIREPO ?= quay.io/lvh-images/lvh
DOCKER ?= docker

all: tests little-vm-helper

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	CGO_ENABLED=0 $(GO) build ./cmd/lvh

.PHONY: image
image:
	$(DOCKER) build -f Dockerfile -t $(OCIREPO) .

.PHONY: install
install:
	CGO_ENABLED=0 $(GO) install ./cmd/lvh

clean:
	rm -f lvh
FORCE:
