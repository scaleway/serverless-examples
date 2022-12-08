FROM golang:1.19-alpine

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY *.html ./

# Build
RUN go build -o /server-image

HEALTHCHECK --interval=10s --timeout=10s --start-period=5s CMD curl -f http://127.0.0.1:8080/status || exit 1

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
#ENV HTTP_PORT=8081

# Run
CMD [ "/server-image" ]
