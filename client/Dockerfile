FROM golang:1.21-alpine3.17 AS builder

WORKDIR /app/
COPY . .

RUN go mod download
RUN set -x; apk add --no-cache && CGO_ENABLED=0 go build -ldflags="-s -w" -o ./client ./client

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/client/index.html .
COPY --from=builder /app/client .

CMD ["./client"]
