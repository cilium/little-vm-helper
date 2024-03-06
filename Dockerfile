FROM --platform=$BUILDPLATFORM golang:1.22.0@sha256:7b297d9abee021bab9046e492506b3c2da8a3722cbf301653186545ecc1e00bb as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:acdc29f25f9c5d678b264007f6f0ea63f12756dcef1cdb043f3fd2a13ba735c3
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
