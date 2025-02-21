FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -o server ./cmd/server

FROM alpine:3.21

COPY --from=builder /app/server /app/server

CMD ["/app/server"]
