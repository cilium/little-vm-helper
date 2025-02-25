FROM --platform=$BUILDPLATFORM golang:1.24.0@sha256:5255fad61a7e8880e742ee3e30ac54d3fdc48ea5236d0bcf14bfedb6643cbeae as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:498a000f370d8c37927118ed80afe8adc38d1edcbfc071627d17b25c88efcab0
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
