FROM --platform=$BUILDPLATFORM golang:1.24.6@sha256:034848561f95a942e2163d9017e672f0c65403f699336db4529a908af00dfc98 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f9a104fddb33220ec80fc45a4e606c74aadf1ef7a3832eb0b05be9e90cd61f5f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
