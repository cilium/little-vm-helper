FROM --platform=$BUILDPLATFORM golang:1.26.2@sha256:2a2b4b5791cea8ae09caecba7bad0bd9631def96e5fe362e4a5e67009fe4ae61 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:1487d0af5f52b4ba31c7e465126ee2123fe3f2305d638e7827681e7cf6c83d5e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
