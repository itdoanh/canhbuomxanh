FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/canhbuomxanh-api ./cmd/api

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
USER appuser
WORKDIR /home/appuser

COPY --from=builder /bin/canhbuomxanh-api ./canhbuomxanh-api

EXPOSE 8080
ENTRYPOINT ["./canhbuomxanh-api"]
