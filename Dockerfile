FROM --platform=$BUILDPLATFORM golang:1.25.1@sha256:ab1f5c47de0f2693ed97c46a646bde2e4f380e40c173454d00352940a379af60 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:d82f458899c9696cb26a7c02d5568f81c8c8223f8661bb2a7988b269c8b9051e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
