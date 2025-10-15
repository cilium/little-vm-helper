FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:6ea52a02734dd15e943286b048278da1e04eca196a564578d718c7720433dbbe as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2f590fc602ce325cbff2ccfc39499014d039546dc400ef8bbf5c6ffb860632e7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
