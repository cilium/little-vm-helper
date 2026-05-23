FROM --platform=$BUILDPLATFORM golang:1.26.3@sha256:6df14f4a4bc9d979a3721f488981e0d1b318006377e473ed23d026796f5f4c0a AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:32015ee641bfecc97161986c9d24957068175444f66fbcbe08a664b6cf5c1c2e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
