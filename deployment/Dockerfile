FROM golang:1.25-alpine AS builder
ARG VERSION=1.0.0
ARG BUILD=1
ARG DATE=04.04.2025
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN echo "Version: $VERSION, Build: $BUILD, Date: $DATE"
RUN export CGO_ENABLED=0
RUN export GOOS=linux
RUN go build -ldflags="-X main.version=$VERSION -X main.build=$BUILD -X main.date=$DATE" -o ./tick ./cmd/main.go
FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/tick /tick/tick
COPY --from=builder /build/web /tick/web

WORKDIR /tick
ENTRYPOINT ["./tick"]