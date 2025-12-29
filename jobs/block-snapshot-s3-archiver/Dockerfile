# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy required files
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

# Build the executable
RUN CGO_ENABLED=0 GOOS=linux go build -o /jobs-snapshot-s3

# Final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /jobs-snapshot-s3 /app/jobs-snapshot-s3

# Run the executable
ENTRYPOINT [ "/app/jobs-snapshot-s3" ]
