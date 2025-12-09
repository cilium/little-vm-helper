FROM --platform=$BUILDPLATFORM golang:1.25.5@sha256:0ece421d4bb2525b7c0b4cad5791d52be38edf4807582407525ca353a429eccc AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:d80cd694d3e9467884fcb94b8ca1e20437d8a501096cdf367a5a1918a34fc2fd
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
