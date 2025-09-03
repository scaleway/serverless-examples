# Using apline/golang image
FROM golang:1.25-alpine

# Set destination for COPY
WORKDIR /app

# Copy required files
COPY go.mod ./
COPY go.sum ./
COPY main.go ./

# Build the executable
RUN go build -o /jobs-snapshot-s3

# Run the executable
ENTRYPOINT [ "/jobs-snapshot-s3" ]
