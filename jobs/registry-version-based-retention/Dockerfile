# Use the alpine version of the golang image as the base image
FROM golang:1.24-alpine

# Set the working directory inside the container to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod ./
COPY go.sum ./

# Copy the Go source files to the working directory
COPY *.go ./

# Build the executable named reg-clean from the Go source files
RUN go build -o /reg-clean

# Set the default command to run the reg-clean executable when the container starts
CMD ["/reg-clean"]
