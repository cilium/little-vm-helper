FROM golang:1.21.5@sha256:672a2286da3ee7a854c3e0a56e0838918d0dbb1c18652992930293312de898a6 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
RUN make little-vm-helper

FROM busybox@sha256:ba76950ac9eaa407512c9d859cea48114eeff8a6f12ebaa5d32ce79d4a017dd8
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
