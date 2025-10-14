FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:7d73c4c57127279b23f3f70cbb368bf0fe08f7ab32af5daf5764173d25e78b74 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:d82f458899c9696cb26a7c02d5568f81c8c8223f8661bb2a7988b269c8b9051e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
