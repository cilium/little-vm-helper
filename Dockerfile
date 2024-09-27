FROM --platform=$BUILDPLATFORM golang:1.23.1@sha256:68c07ac27294fbdee3130831a4a6af66b9720b3cedd854c07e87e1625d95e848 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:c230832bd3b0be59a6c47ed64294f9ce71e91b327957920b6929a0caa8353140
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
