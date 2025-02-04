FROM --platform=$BUILDPLATFORM golang:1.23.5@sha256:e213430692e5c31aba27473cdc84cfff2896d0c097e984bef67b6a44c75a8181 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:a5d0ce49aa801d475da48f8cb163c354ab95cab073cd3c138bd458fc8257fbf1
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
