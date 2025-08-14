FROM --platform=$BUILDPLATFORM golang:1.24.6@sha256:e155b5162f701b7ab2e6e7ea51cec1e5f6deffb9ab1b295cf7a697e81069b050 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f9a104fddb33220ec80fc45a4e606c74aadf1ef7a3832eb0b05be9e90cd61f5f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
