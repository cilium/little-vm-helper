FROM --platform=$BUILDPLATFORM golang:1.22.2@sha256:c54c7d60b3bf561264bd6ef1f88cb6c4b5ec60d5881caf998e75d3699a0d098b as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:6776a33c72b3af7582a5b301e3a08186f2c21a3409f0d2b52dfddbdbe24a5b04
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
