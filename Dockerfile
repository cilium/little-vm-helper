FROM --platform=$BUILDPLATFORM golang:1.24.3@sha256:4c0a1814a7c6c65ece28b3bfea14ee3cf83b5e80b81418453f0e9d5255a5d7b8 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
