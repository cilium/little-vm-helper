FROM --platform=$BUILDPLATFORM golang:1.23.2@sha256:858ab89651d8a3d637da5580e71fdec40b5aefbb148ba50b9a629bd079a14bcd as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:768e5c6f5cb6db0794eec98dc7a967f40631746c32232b78a3105fb946f3ab83
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
