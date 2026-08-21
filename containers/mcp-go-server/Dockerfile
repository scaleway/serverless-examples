# syntax=docker/dockerfile:1

# Step 1: Build binary using modern Go toolchain
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

# BuildKit cache for module downloads
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY main.go .

# Build with BuildKit compiler cache + binary size optimizations
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-w -s" \
    -o mcp-server main.go

# Step 2: Minimal non-root distroless runtime (~5MB final image)
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

COPY --from=builder /app/mcp-server /mcp-server

EXPOSE 8080
ENV PORT=8080

USER nonroot:nonroot

ENTRYPOINT ["/mcp-server"]
