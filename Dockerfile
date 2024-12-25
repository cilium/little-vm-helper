FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:f06d2bb355a67ccc6c23f3699766323a09ed0a4b724a6b25a300d4b30e01f02c as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2919d0172f7524b2d8df9e50066a682669e6d170ac0f6a49676d54358fe970b5
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
