GO ?= go

OCIREPO ?= quay.io/lvh-images/lvh
DOCKER ?= docker
VERSION ?= $(shell git describe --tags --always --long)

GO_BUILD_LDFLAGS =
GO_BUILD_LDFLAGS += -X 'github.com/cilium/little-vm-helper/pkg/version.Version=$(VERSION)'
GO_BUILD_FLAGS += -ldflags "$(GO_BUILD_LDFLAGS)"

UNAME_M := $(shell uname -m)
ifeq ($(UNAME_M),x86_64)
	TARGET_ARCH ?= amd64
else ifeq ($(UNAME_M),aarch64)
	TARGET_ARCH ?= arm64
else
	TARGET_ARCH ?= amd64
endif

all: tests little-vm-helper

.PHONY: tests
tests:
	$(GO) test -cover ./...

little-vm-helper: FORCE
	GOARCH=${TARGET_ARCH} CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) ./cmd/lvh

.PHONY: image
image:
	$(DOCKER) build -f Dockerfile --platform=linux/${TARGET_ARCH} -t $(OCIREPO) .

.PHONY: install
install:
	GOARCH=${TARGET_ARCH} CGO_ENABLED=0 $(GO) install ./cmd/lvh

clean:
	rm -f lvh
FORCE:
