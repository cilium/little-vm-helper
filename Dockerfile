FROM --platform=$BUILDPLATFORM golang:1.26.5@sha256:6cd10a6fcc5eadd62008fc2ad8056b38971cafd42f44d55297f18be8adc86410 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:dc2d74b28e4cf8984fa52af1f39bc7c3d9c73760b41a74d629f5d11b1ab28616
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
