FROM --platform=$BUILDPLATFORM golang:1.25.4@sha256:e68f6a00e88586577fafa4d9cefad1349c2be70d21244321321c407474ff9bf2 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:e3652a00a2fabd16ce889f0aa32c38eec347b997e73bd09e69c962ec7f8732ee
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
