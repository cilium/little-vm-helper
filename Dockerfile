FROM golang:1.21.6@sha256:76aadd914a29a2ee7a6b0f3389bb2fdb87727291d688e1d972abe6c0fa6f2ee0 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
RUN make little-vm-helper

FROM busybox@sha256:ba76950ac9eaa407512c9d859cea48114eeff8a6f12ebaa5d32ce79d4a017dd8
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
