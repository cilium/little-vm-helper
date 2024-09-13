FROM --platform=$BUILDPLATFORM golang:1.23.1@sha256:2fe82a3f3e006b4f2a316c6a21f62b66e1330ae211d039bb8d1128e12ed57bf1 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:c230832bd3b0be59a6c47ed64294f9ce71e91b327957920b6929a0caa8353140
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
