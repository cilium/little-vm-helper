GO ?= go

OCIREPO ?= quay.io/lvh-images/lvh
DOCKER ?= docker

all: tests little-vm-helper 

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	$(GO) build ./cmd/lvh

.PHONY: image
image:
	$(DOCKER) build -f Dockerfile -t $(OCIREPO) .

clean:
	rm -f lvh
FORCE:
