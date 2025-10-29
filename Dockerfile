FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:6bac879c5b77e0fc9c556a5ed8920e89dab1709bd510a854903509c828f67f96 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:fba0711bd6995f7e0158f397d85f63998b4c8b1a1e3b1e9e0394c7b165585440
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
