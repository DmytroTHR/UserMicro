FROM golang:1.17-alpine3.13 as builder
WORKDIR /go/src/UserMicro
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o build/userservice

FROM scratch
COPY --from=builder /go/src/UserMicro/build/userservice /usr/bin/userservice
ENTRYPOINT [ "/usr/bin/userservice" ]