FROM --platform=$BUILDPLATFORM golang:1.23.1@sha256:8f6a7d881c5348114a19a08f8cc052b3a0d5341539e7ecc9e00902776949bc71 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:c230832bd3b0be59a6c47ed64294f9ce71e91b327957920b6929a0caa8353140
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
