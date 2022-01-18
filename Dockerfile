FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/UserMicro
COPY . .
RUN apk add --no-cache make
RUN make go-build

FROM scratch
COPY --from=builder /go/src/UserMicro/bin/userservice /usr/bin/userservice
ENTRYPOINT [ "/usr/bin/userservice" ]