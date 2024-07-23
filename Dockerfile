FROM --platform=$BUILDPLATFORM golang:1.22.5@sha256:1b70fa8200ed48367e7d0e0e82c0c1d14cf712d422116f9d69e7c964243653a0 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
