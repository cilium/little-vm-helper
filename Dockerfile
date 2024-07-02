FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:800fde66644c2447abc2a085cbecdcb53d2e21234ec8e80be0c90fc011a93f4b as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
