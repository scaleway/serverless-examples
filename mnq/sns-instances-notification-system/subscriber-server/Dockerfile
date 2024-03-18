FROM debian:bookworm-20240211

RUN apt-get update && \
    apt-get install -y golang ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./www ./www

RUN go mod init subscriber-server && \
    go mod tidy && \
    go build -o ./subscriber ./cmd/subscriber

EXPOSE 8081

CMD ["./subscriber"]