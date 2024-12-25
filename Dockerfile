FROM --platform=$BUILDPLATFORM golang:1.23.4@sha256:cb00c135e0e60590614a68f1361fcafdf0f4d2bec2efa7d11b4eb634d1cbe4ed as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:2919d0172f7524b2d8df9e50066a682669e6d170ac0f6a49676d54358fe970b5
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
