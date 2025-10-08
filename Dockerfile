FROM --platform=$BUILDPLATFORM golang:1.25.2@sha256:1c91b4f4391774a73d6489576878ad3ff3161ebc8c78466ec26e83474855bfcf as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:d82f458899c9696cb26a7c02d5568f81c8c8223f8661bb2a7988b269c8b9051e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
