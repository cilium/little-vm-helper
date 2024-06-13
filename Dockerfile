FROM --platform=$BUILDPLATFORM golang:1.22.4@sha256:245c739b7dc876b87e90eabfdbab6bfb14e7e909b7bf260013594937521a311c as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
