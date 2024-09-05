FROM --platform=$BUILDPLATFORM golang:1.23.0@sha256:acfb46be39840f8c2a6b9efdd673c6627011200c73bab4e6d18b8b9ab4641c46 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:34b191d63fbc93e25e275bfccf1b5365664e5ac28f06d974e8d50090fbb49f41
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
