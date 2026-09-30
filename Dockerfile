FROM --platform=$BUILDPLATFORM golang:1.27.1@sha256:e0174e51e81218523251d85d248a90d24c3d5e81543b4f07a5d66229397db190 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
