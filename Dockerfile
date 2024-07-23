FROM --platform=$BUILDPLATFORM golang:1.22.5@sha256:695e2559491efb2cc266226501b128eb6b4923d388f55ec182e1d96f65955a2a as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:9ae97d36d26566ff84e8893c64a6dc4fe8ca6d1144bf5b87b2b85a32def253c7
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
