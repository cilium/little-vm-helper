FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:0afe9b553123ec160b8282154b7156241ebf43c9033ea2bf9528c6c58c17e89d AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:e3652a00a2fabd16ce889f0aa32c38eec347b997e73bd09e69c962ec7f8732ee
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
