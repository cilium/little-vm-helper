FROM --platform=$BUILDPLATFORM golang:1.22.3@sha256:f43c6f049f04cbbaeb28f0aad3eea15274a7d0a7899a617d0037aec48d7ab010 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:82078ce09d47857439e17eb6a35b0c24ad5be8cde11f56f81c4e293029da727b
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
