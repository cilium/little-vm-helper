FROM --platform=$BUILDPLATFORM golang:1.24.5@sha256:267159cb984d1d034fce6e9db8641bf347f80e5f2e913561ea98c40d5051cb67 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
