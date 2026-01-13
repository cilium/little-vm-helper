FROM --platform=$BUILDPLATFORM golang:1.25.5@sha256:581c059c1d53b96a6db51a2f7a0eb943491ba1da50f1e43700db0cc325618096 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2383baad1860bbe9d8a7a843775048fd07d8afe292b94bd876df64a69aae7cb1
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
