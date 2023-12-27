FROM golang:1.21-alpine3.17 AS builder

WORKDIR /app/

# Copy go mod and sum files
COPY . .

RUN go mod download
RUN set -x; apk add --no-cache && CGO_ENABLED=0 go build -ldflags="-s -w" -o ./server  ./server

# Use alpine for the final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]
