FROM golang:1.24.1-bookworm

# Set destination for COPY
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

# Build
RUN go build -o /server-image

# Run
CMD [ "/server-image" ]
