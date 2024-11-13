FROM --platform=$BUILDPLATFORM golang:1.23.3@sha256:26602182e32576fca9ad5397fd74b87ae68f7ee727436a85dd19a1c3766ad685 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:768e5c6f5cb6db0794eec98dc7a967f40631746c32232b78a3105fb946f3ab83
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
