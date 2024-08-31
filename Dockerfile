FROM --platform=$BUILDPLATFORM golang:1.23.0@sha256:613a108a4a4b1dfb6923305db791a19d088f77632317cfc3446825c54fb862cd as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:82742949a3709938cbeb9cec79f5eaf3e48b255389f2dcedf2de29ef96fd841c
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
