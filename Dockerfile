FROM --platform=$BUILDPLATFORM golang:1.24.0@sha256:58cf31c1b10858d5772475e4f86853a03e0e220cc6cfd60965053951a2563b5e as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:498a000f370d8c37927118ed80afe8adc38d1edcbfc071627d17b25c88efcab0
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
