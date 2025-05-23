FROM --platform=$BUILDPLATFORM golang:1.24.3@sha256:795a40cbe36a11e39b0709eb98d7aaa8d312d60336d863a0ef1c8aff07c1b3e0 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f64ff79725d0070955b368a4ef8dc729bd8f3d8667823904adcb299fe58fc3da
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
