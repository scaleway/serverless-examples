FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    openssl \
    wget

ARG PGP_URL=https://www.mongodb.org/static/pgp
ARG REPO_URL=https://repo.mongodb.org/apt/ubuntu

ARG MONGODB_VERSION=7.0
ARG MONGOSH_VERSION=2.1.1

RUN wget -qO- ${PGP_URL}/server-${MONGODB_VERSION}.asc \
    | tee /etc/apt/trusted.gpg.d/server-${MONGODB_VERSION}.asc

RUN echo "deb ${REPO_URL} jammy/mongodb-org/${MONGODB_VERSION} multiverse" \
    | tee /etc/apt/sources.list.d/mongodb-org-${MONGODB_VERSION}.list

RUN apt-get update
RUN apt-get install -y mongodb-mongosh

WORKDIR /app

COPY mongo-check.sh .
CMD [ "./mongo-check.sh" ]
