FROM --platform=$BUILDPLATFORM golang:1.27.2@sha256:007a316059246276cb214df93950be77837545dc78007a73e0291a4963d5d3bf AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:fd7dc98638c8e305f4dc34e979f1c0fdfdcaeb0fbf8fcff77ae834b6da3d7e6e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
