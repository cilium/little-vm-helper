FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:b01f7c744a3f1fccaf44905169169fed0ab13e6d1d702a6542d07b34cf677969 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2919d0172f7524b2d8df9e50066a682669e6d170ac0f6a49676d54358fe970b5
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
