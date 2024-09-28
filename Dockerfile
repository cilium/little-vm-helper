FROM --platform=$BUILDPLATFORM golang:1.23.1@sha256:efa59042e5f808134d279113530cf419e939d40dab6475584a13c62aa8497c64 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:c230832bd3b0be59a6c47ed64294f9ce71e91b327957920b6929a0caa8353140
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
