FROM --platform=$BUILDPLATFORM golang:1.25.1@sha256:1fd7d46f956287d1856b92add5cc5ab8b87c07a1ed766419bb603a8620746b4a as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
