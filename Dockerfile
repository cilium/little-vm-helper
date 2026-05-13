FROM --platform=$BUILDPLATFORM golang:1.26.3@sha256:633d23bf362cb40dd72b4f277288a8929697d77537f9c801b81aeced19b5bdf3 AS gobuilder
WORKDIR /src/little-vm-helper
COPY . .
ARG TARGETARCH
RUN TARGET_ARCH=$TARGETARCH make little-vm-helper

FROM busybox@sha256:1487d0af5f52b4ba31c7e465126ee2123fe3f2305d638e7827681e7cf6c83d5e
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
