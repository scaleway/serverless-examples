# syntax=docker/dockerfile:1

FROM node:20-alpine

WORKDIR /usr/app
COPY ./ ./
RUN npm install
RUN npm run build

#Web application configuration
ENV PORT=8080

#Database configuration used for dynamically rendered data
ENV PGHOST=localhost
ENV PGPORT=5432
ENV PGDATABASE=database
ENV PGUSER=user
ENV PGPASSWORD=password

CMD npm run start