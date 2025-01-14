FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:95983c2237525489841f6c9dc8e0df3d5cfbcf7788d74a6720ef5163ab7fc0f7 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:a5d0ce49aa801d475da48f8cb163c354ab95cab073cd3c138bd458fc8257fbf1
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
