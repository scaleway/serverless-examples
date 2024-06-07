# Use the official Rust image as the build stage
FROM rust:1.78.0-alpine3.20 as builder

# Set the working directory inside the container
WORKDIR /usr/src/app

# Copy the source code to the working directory
COPY . .

# Build the Rust application
RUN cargo build --release

# Use a smaller base image for the final container
FROM alpine:3.20

# Copy the compiled binary from the build stage
COPY --from=builder /usr/src/app/target/release/hello-world /usr/local/bin/hello-world

# Expose the port that your application will run on
EXPOSE 8080

# Set the entrypoint command
CMD ["hello-world"]
