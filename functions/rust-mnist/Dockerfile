FROM rust:1.65-alpine3.15 as builder

RUN apk upgrade --update-cache --available && \
    apk add musl-dev

ENV HOME=/home/root
WORKDIR $HOME/rust-mnist

COPY Cargo.toml ./
COPY Cargo.lock ./

COPY common ./common
COPY function ./function
COPY training ./training

RUN cargo build --bin test-server --release && \
    # Copy executable out of the cache so it is available in the final image.
    cp target/release/test-server ./test-server

FROM alpine:3.15

COPY --from=builder /home/root/rust-mnist/test-server /test-server

EXPOSE 3000
CMD [ "/test-server" ]