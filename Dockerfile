FROM --platform=$BUILDPLATFORM golang:1.22.6@sha256:305aae5c6d68e9c122246147f2fde17700d477c9062b6724c4d182abf6d3907e as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
