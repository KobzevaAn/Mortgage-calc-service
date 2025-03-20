FROM golang:1.22.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o mortgage-calc-service ./cmd

RUN go mod vendor

FROM alpine:latest

WORKDIR /root/

COPY configs/config.yaml ./configs/config.yaml
COPY --from=builder /app/mortgage-calc-service .

CMD ["./mortgage-calc-service"]
