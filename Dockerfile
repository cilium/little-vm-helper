FROM golang:1.21.7@sha256:8144f2d44d2262fa930b437200fc4ada624d8a0b9c83d688e2a6f545d097c45b as gobuilder
WORKDIR /src/little-vm-helper
COPY . .
RUN make little-vm-helper

FROM busybox@sha256:6d9ac9237a84afe1516540f40a0fafdc86859b2141954b4d643af7066d598b74
COPY --from=gobuilder /src/little-vm-helper/lvh /usr/bin/lvh
