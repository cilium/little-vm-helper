FROM --platform=$BUILDPLATFORM golang:1.24.3@sha256:795a40cbe36a11e39b0709eb98d7aaa8d312d60336d863a0ef1c8aff07c1b3e0 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:37f7b378a29ceb4c551b1b5582e27747b855bbfaa73fa11914fe0df028dc581f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
