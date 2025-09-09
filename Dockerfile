FROM --platform=$BUILDPLATFORM golang:1.25.1@sha256:d6bdb04c3109d29739c19566092e63b469d4cdf193df206fadf286aa4a549ba6 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:ab33eacc8251e3807b85bb6dba570e4698c3998eca6f0fc2ccb60575a563ea74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
