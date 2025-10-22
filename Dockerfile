FROM --platform=$BUILDPLATFORM golang:1.25.3@sha256:8c945d3e25320e771326dafc6fb72ecae5f87b0f29328cbbd87c4dff506c9135 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2f590fc602ce325cbff2ccfc39499014d039546dc400ef8bbf5c6ffb860632e7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
