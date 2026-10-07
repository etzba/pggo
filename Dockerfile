# build
FROM golang:1.24-bookworm AS builder
LABEL authors="support@etzba.com, Nadav Ben Mazia"
COPY . /build
WORKDIR /build
# Go mod download and verify dependencies
RUN go mod download
RUN go mod verify
# Build the app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o pggo main.go

# New distroless image with no root
FROM gcr.io/distroless/static:nonroot
# Copy the app from builder
COPY --from=builder /build/pggo /pggo
WORKDIR /
USER 65532:65532
CMD ["./pggo"]

