FROM golang:1.23-alpine

# Set destination for COPY
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

# Build
RUN go build -o /server-image

# Run
CMD [ "/server-image" ]
