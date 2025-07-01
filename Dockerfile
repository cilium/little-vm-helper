FROM --platform=$BUILDPLATFORM golang:1.24.4@sha256:270cd5365c84dd24716c42d7f9f7ddfbc131c8687e163e6748b9c1322c518213 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
