FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:9820aca42262f58451f006de3213055974b36f24b31508c1baa73c967fcecb99 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:a5d0ce49aa801d475da48f8cb163c354ab95cab073cd3c138bd458fc8257fbf1
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
