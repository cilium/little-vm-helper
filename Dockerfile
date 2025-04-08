FROM --platform=$BUILDPLATFORM golang:1.24.2@sha256:b51b7beeabe2e2d8438ba4295c59d584049873a480ba0e7b56d80db74b3e3a3a as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:37f7b378a29ceb4c551b1b5582e27747b855bbfaa73fa11914fe0df028dc581f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
