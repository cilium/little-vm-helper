FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:3589439790974ec05491b66b656bf1048d0f50dd010a903463e3156ba1fc26de as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
