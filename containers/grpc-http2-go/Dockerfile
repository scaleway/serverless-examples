FROM golang:1.22-bookworm as build

# Set working directory
WORKDIR /app

# Copy required files
COPY go.mod .
COPY go.sum .
COPY server/*.go ./server/
COPY hello.proto ./

# We install the protobuf compilation stack on the image directly to simplify workflow
RUN apt-get update && apt-get install -y protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Generate go protobuf files
RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative hello.proto

# Build the executable
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build server/main.go

# Put the executable in a light scratch image
FROM scratch

COPY --from=build /app/main /server

ENTRYPOINT ["/server"]