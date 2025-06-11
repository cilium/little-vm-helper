FROM --platform=$BUILDPLATFORM golang:1.24.4@sha256:d1db785fb37feb87d9140d78ed4fb7c75ee787360366a9c5efe39c7a841a0277 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f85340bf132ae937d2c2a763b8335c9bab35d6e8293f70f606b9c6178d84f42b
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
