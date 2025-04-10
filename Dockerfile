FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o grpc-server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/grpc-server .

EXPOSE 50051

CMD ["./grpc-server"]