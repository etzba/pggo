# build
FROM golang:1.22-bullseye AS be_builder

COPY . /build

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pggo main.go

# alpine
FROM alpine:3.20

RUN apk add ca-certificates

COPY --from=be_builder /build/pggo /pggo

WORKDIR /

CMD ["./pggo"]