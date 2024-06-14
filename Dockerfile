FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:a0679accac8685cc5389bd2298e045e570100940e6bdcca666a8ca7b32a1276c as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
