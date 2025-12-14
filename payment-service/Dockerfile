FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment_server ./cmd/server/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/payment_server .
EXPOSE 50051
CMD ["./payment_server"]