FROM --platform=$BUILDPLATFORM golang:1.25.1@sha256:76a94c4a37aaab9b1b35802af597376b8588dc54cd198f8249633b4e117d9fcc as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
