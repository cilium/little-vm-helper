FROM --platform=$BUILDPLATFORM golang:1.23.3@sha256:73f06be4578c9987ce560087e2e2ea6485fb605e3910542cadd8fa09fc5f3e31 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:5b0f33c83a97f5f7d12698df6732098b0cdb860d377f6307b68efe2c6821296f
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
