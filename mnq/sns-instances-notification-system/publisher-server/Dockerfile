FROM debian:bookworm-20240211
LABEL authors="jlecomte@scaleway.com"

RUN apt-get update && \
    apt-get install -y golang ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./www ./www

RUN go mod init publisher-server && \
    go mod tidy

EXPOSE 8081

WORKDIR /app/cmd/publisher

CMD ["go", "run", "main.go"]