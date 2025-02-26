FROM --platform=$BUILDPLATFORM golang:1.24.0@sha256:cd0c949a4709ef70a8dad14274f09bd07b25542de5a1c4812f217087737efd17 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:498a000f370d8c37927118ed80afe8adc38d1edcbfc071627d17b25c88efcab0
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
