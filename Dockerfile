FROM --platform=$BUILDPLATFORM golang:1.24.2@sha256:1ecc479bc712a6bdb56df3e346e33edcc141f469f82840bab9f4bc2bc41bf91d as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:37f7b378a29ceb4c551b1b5582e27747b855bbfaa73fa11914fe0df028dc581f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
