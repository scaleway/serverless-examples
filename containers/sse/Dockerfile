FROM golang:1.22-alpine3.19 AS build

# Install git, this used to fetch VCS information when building the application
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/sse

FROM alpine:3.19

WORKDIR /

COPY ./www /www

COPY --from=build /app/sse .

EXPOSE 8080

CMD ["/sse"]
