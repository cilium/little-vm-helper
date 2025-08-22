FROM --platform=$BUILDPLATFORM golang:1.25.0@sha256:4859242e2c392ddc9d3225fd41181c00a443d9cc005b8e5131ce164106fbc676 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
