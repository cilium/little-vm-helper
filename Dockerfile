FROM --platform=$BUILDPLATFORM golang:1.25.0@sha256:5502b0e56fca23feba76dbc5387ba59c593c02ccc2f0f7355871ea9a0852cebe as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
