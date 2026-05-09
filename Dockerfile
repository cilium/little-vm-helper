FROM --platform=$BUILDPLATFORM golang:1.26.3@sha256:2981696eed011d747340d7252620932677929cce7d2d539602f56a8d7e9b660b AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:1487d0af5f52b4ba31c7e465126ee2123fe3f2305d638e7827681e7cf6c83d5e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
