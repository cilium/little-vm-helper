FROM --platform=$BUILDPLATFORM golang:1.25.6@sha256:ce63a16e0f7063787ebb4eb28e72d477b00b4726f79874b3205a965ffd797ab2 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:b86c79d6d337dcaa4cad3cfc704d5e6ae5ba725f3e0674d96f65af7a8f5e2761
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
