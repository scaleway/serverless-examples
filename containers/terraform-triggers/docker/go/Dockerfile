FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY . ./
RUN go build -o server

FROM alpine

COPY --from=builder /app/server /app/server

CMD ["/app/server"]
