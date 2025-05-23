FROM --platform=$BUILDPLATFORM golang:1.24.3@sha256:4c0a1814a7c6c65ece28b3bfea14ee3cf83b5e80b81418453f0e9d5255a5d7b8 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:f64ff79725d0070955b368a4ef8dc729bd8f3d8667823904adcb299fe58fc3da
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
