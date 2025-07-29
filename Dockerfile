FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o mini-alt .

FROM debian:bookworm-slim

ENV XDG_CONFIG_HOME=/data

RUN mkdir -p /data

WORKDIR /app

COPY --from=builder /app/mini-alt .

EXPOSE 9000
EXPOSE 9001

CMD ["./mini-alt"]
