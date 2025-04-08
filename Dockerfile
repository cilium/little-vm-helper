FROM --platform=$BUILDPLATFORM golang:1.24.2@sha256:fb224f950b0dfb889f8f1122bf3dfc21976735ccf983b73a1e6e014215a272f5 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:37f7b378a29ceb4c551b1b5582e27747b855bbfaa73fa11914fe0df028dc581f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
