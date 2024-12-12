FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:70031844b8c225351d0bb63e2c383f80db85d92ba894e3da7e13bcf80efa9a37 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2919d0172f7524b2d8df9e50066a682669e6d170ac0f6a49676d54358fe970b5
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
