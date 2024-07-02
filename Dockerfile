FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:3e731b22f6db4e5a76d957250d9a2bf9c2e2717ff7f97cfd20322e81ba185898 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
