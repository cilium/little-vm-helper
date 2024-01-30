GO ?= go

OCIREPO ?= quay.io/lvh-images/lvh
DOCKER ?= docker
VERSION ?= $(shell git describe --tags --always --long)

GO_BUILD_LDFLAGS =
GO_BUILD_LDFLAGS += -X 'github.com/cilium/little-vm-helper/pkg/version.Version=$(VERSION)'
GO_BUILD_FLAGS += -ldflags "$(GO_BUILD_LDFLAGS)"


all: tests little-vm-helper

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) ./cmd/lvh

.PHONY: image
image:
	$(DOCKER) build -f Dockerfile -t $(OCIREPO) .

.PHONY: install
install:
	CGO_ENABLED=0 $(GO) install ./cmd/lvh

clean:
	rm -f lvh
FORCE:
