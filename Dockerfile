FROM --platform=$BUILDPLATFORM golang:1.22.6@sha256:d5e49f92b9566b0ddfc59a0d9d85cd8a848e88c8dc40d97e29f306f07c3f8338 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
