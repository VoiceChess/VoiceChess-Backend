FROM --platform=linux/amd64 ubuntu:22.04 AS builder

RUN apt-get update && apt-get install -y wget gcc libc6-dev && \
    wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz && \
    rm go1.24.0.linux-amd64.tar.gz && \
    rm -rf /var/lib/apt/lists/*

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -o main .

FROM --platform=linux/amd64 ubuntu:22.04

ENV PATH="/usr/games:${PATH}"

RUN apt-get update && apt-get install -y stockfish && rm -rf /var/lib/apt/lists/*
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /root/
COPY .env.docker .env
COPY --from=builder /app/main .
COPY --from=builder /app/schema ./schema

EXPOSE 8080
CMD ["./main"]