FROM --platform=$BUILDPLATFORM golang:1.23.3@sha256:e5ca1999e21764b1fd40cf6564ebfb7022e7a55b8c72886a9bcb697a5feac8d6 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:db142d433cdde11f10ae479dbf92f3b13d693fd1c91053da9979728cceb1dc68
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
