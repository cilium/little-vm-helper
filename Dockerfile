FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:dd08f769578a5f51a22bf6a81109288e23cfe2211f051a5c29bd1c05ad3db52a AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2f590fc602ce325cbff2ccfc39499014d039546dc400ef8bbf5c6ffb860632e7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
