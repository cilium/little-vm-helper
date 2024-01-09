FROM golang:1.21.5 as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
RUN make little-vm-helper

FROM busybox
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
