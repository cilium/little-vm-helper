FROM --platform=$BUILDPLATFORM golang:1.27.1@sha256:f773aa2679c165b2d4ccf04c1de9ef5f34c0e4fd7008b3c449d4cdc47d83af55 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:dc2d74b28e4cf8984fa52af1f39bc7c3d9c73760b41a74d629f5d11b1ab28616
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
