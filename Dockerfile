FROM --platform=$BUILDPLATFORM golang:1.25.1@sha256:d7098379b7da665ab25b99795465ec320b1ca9d4addb9f77409c4827dc904211 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:d82f458899c9696cb26a7c02d5568f81c8c8223f8661bb2a7988b269c8b9051e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
