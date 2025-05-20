# Use the official Rust image as the build stage
FROM rust:1.87.0-alpine3.21 as builder

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the source code to the working directory
COPY . .

# Build the Rust application
RUN cargo build --release

# Use a smaller base image for the final container
FROM alpine:3.21

# Copy the compiled binary from the build stage
COPY --from=builder /usr/src/app/target/release/hello-world /usr/local/bin/hello-world

# Set the entrypoint command
CMD ["hello-world"]
